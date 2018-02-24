// Copyright 2009 smallnest. All rights reserved.
// Use of this source code is governed by Apache License Version 2.0
// license that can be found in the LICENSE file.

package glean

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"plugin"
	"reflect"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/smallnest/glean/log"
	fsnotify "gopkg.in/fsnotify.v1"
)

var (
	// ErrItemHasNotConfigured the plugin item has not been configured.
	ErrItemHasNotConfigured = errors.New("pluginItem is not configured")
	// ErrClosed glean instance has been closed.
	ErrClosed = errors.New("glean has been closed")
)

// PluginItem is a configured item that can be reloaded.
type PluginItem struct {
	// File file path of this plugin.
	File string `json:"file"`
	// ID is an unique string for this item.
	ID string `json:"id"`
	// Name is name of the symbol. Notice id is unique but names may be duplicated in different plugins.
	Name string `json:"name"`
	// Version is version of the plugin for tracing and upgrade.
	Version string `json:"version"`
	// Cached points the opened plugin.
	Cached *plugin.Plugin `json:"-"`
	// v is the function or variable that can be reloaded.
	v interface{}
}

// Glean is a manager that manages all configured plugins and reloaded objects.
type Glean struct {
	configFile  string
	pluginItems []*PluginItem
	idMap       map[string]*PluginItem
	watched     map[string]bool
	mu          sync.RWMutex
	done        chan bool
	closed      bool
}

// New returns a new Glean.
func New(configFile string) *Glean {
	return &Glean{
		configFile: configFile,
		watched:    make(map[string]bool),
		idMap:      make(map[string]*PluginItem),
		done:       make(chan bool),
	}
}

// Close closes Glean and stop watching.
func (g *Glean) Close() {
	g.mu.Lock()
	if !g.closed {
		g.closed = true
		close(g.done)
		g.pluginItems = []*PluginItem{}
		g.idMap = nil
		g.watched = nil
	}
	g.mu.Unlock()
}

// LoadConfig loads plugins from the configured file.
func (g *Glean) LoadConfig() (err error) {
	buf, err := ioutil.ReadFile(g.configFile)
	if err != nil {
		log.Errorf("failed to load %s: %v", g.configFile, err)
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	err = json.Unmarshal(buf, &(g.pluginItems))
	if err != nil {
		log.Errorf("failed to unmarshal %s: %v", g.configFile, err)
		return err
	}

	// initial plugin
	for _, item := range g.pluginItems {
		pp, err := plugin.Open(item.File)
		if err != nil {
			log.Errorf("failed to load %s: %v", item.Name, err)
			return err
		}

		item.Cached = pp
		g.idMap[item.ID] = item
	}

	// watch changes
	err = g.startWatch()
	return
}

// start to watch changes of config changes
func (g *Glean) startWatch() error {
	if g.closed {
		return ErrClosed
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = watcher.Add(g.configFile)
	if err != nil {
		watcher.Close()
		log.Fatal(err)
	}

	go func() {
	watch:
		for {
			select {
			case event := <-watcher.Events:
				log.Info("watch event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Infof("config file %s is modified", event.Name)
					g.checkChanges() // the config file has been modified
				}
			case err := <-watcher.Errors:
				log.Errorf("watcher error: %v", err)
			case <-g.done:
				break watch
			}
		}
	}()

	return err
}

func (g *Glean) checkChanges() {
	buf, err := ioutil.ReadFile(g.configFile)
	if err != nil {
		log.Errorf("failed to load %s: %v", g.configFile, err)
		return
	}

	var latestPluginItems []*PluginItem

	err = json.Unmarshal(buf, &(latestPluginItems))
	if err != nil {
		log.Errorf("failed to unmarshal %s: %v", g.configFile, err)
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	currentPluginItems := g.pluginItems
	added, changed, removed := diffPlugins(currentPluginItems, latestPluginItems)

	for _, item := range removed {
		delete(g.idMap, item.ID)
		delete(g.watched, item.ID)
	}

	// update changed
	for _, item := range changed {
		pp, e := plugin.Open(item.File)
		if e != nil {
			log.Errorf("failed to load %s: %v", item.Name, e)
			err = multierror.Append(err, e)
		}
		item.Cached = pp
		item.v = g.idMap[item.ID].v
		g.idMap[item.ID] = item
	}

	// add added
	for _, item := range added {
		pp, e := plugin.Open(item.File)
		if e != nil {
			log.Errorf("failed to load %s: %v", item.Name, e)
			err = multierror.Append(err, e)
		}
		item.Cached = pp
		g.idMap[item.ID] = item
	}

	//reload all variables
	for _, item := range changed {
		watchID := g.watched[item.ID]

		if watchID {
			e := ReloadFromPlugin(item.Cached, item.Name, item.v)
			if e != nil {
				log.Errorf("failed to reload %s, %s from %s: %v", item.ID, item.Name, item.File, err)
				err = multierror.Append(err, e)
			} else {
				log.Infof("succeeded to reload %s, %s from %s", item.ID, item.Name, item.File)
			}
		}
	}
}

func diffPlugins(currentPluginItems, latestPluginItems []*PluginItem) (added, changed, removed []*PluginItem) {
	latestM := make(map[string]*PluginItem)
	for _, item := range latestPluginItems {
		latestM[item.ID] = item
	}

	currentM := make(map[string]*PluginItem)
	for _, item := range currentPluginItems {
		currentM[item.ID] = item
	}

	for _, item := range latestPluginItems {
		if i, exist := currentM[item.ID]; exist {
			if item.File != i.File {
				changed = append(changed, item)
			}
		} else {
			added = append(added, item)
		}
	}

	for _, item := range currentPluginItems {
		if _, exist := latestM[item.ID]; !exist {
			removed = append(removed, item)
		}
	}

	return
}

// Reload loads an variable or function from configured plugins.
func (g *Glean) Reload(id string, vPtr interface{}) error {
	if g.closed {
		return ErrClosed
	}

	g.mu.RLock()
	item := g.idMap[id]
	g.mu.RUnlock()

	if item == nil {
		return ErrItemHasNotConfigured
	}

	return ReloadFromPlugin(item.Cached, item.Name, vPtr)
}

// Watch watches plugin changes and reload given function/variable automatically.
func (g *Glean) Watch(id string, vPtr interface{}) {
	g.mu.Lock()
	g.watched[id] = true
	item := g.idMap[id]
	g.mu.Unlock()
	if item != nil {
		item.v = vPtr
	}
}

// ReloadAndWatch loads an variable or function from plugins and begin to watch.
func (g *Glean) ReloadAndWatch(id string, vPtr interface{}) error {
	err := g.Reload(id, vPtr)
	if err != nil {
		return err
	}

	g.Watch(id, vPtr)

	return nil
}

// GetObjectByID gets the variable or function by ID.
func (g *Glean) GetObjectByID(id string) (v interface{}) {
	g.mu.RLock()
	v = g.idMap[id].v
	g.mu.RUnlock()
	return v
}

// GetSymbolByID gets the variable or function by ID from cached plugin.
func (g *Glean) GetSymbolByID(id string) (v interface{}, err error) {
	g.mu.RLock()
	v, err = g.idMap[id].Cached.Lookup(g.idMap[id].Name)
	g.mu.RUnlock()
	return v, err
}

// FindAllPlugins gets all IDs that implements interface t.
func (g *Glean) FindAllPlugins(t reflect.Type) ([]string, error) {
	if t.Kind() != reflect.Interface {
		return nil, errors.New("parameter i is not an interface type")
	}

	var ids []string
	g.mu.RLock()
	for id, item := range g.idMap {
		if item.v != nil {
			if reflect.TypeOf(item.v).Implements(t) {
				ids = append(ids, id)
			}
		} else {
			s, err := item.Cached.Lookup(item.Name)
			if err == nil && reflect.ValueOf(s).Elem().Type().Implements(t) {
				ids = append(ids, id)
			}
		}
	}

	return ids, nil
}
