// Copyright 2009 smallnest. All rights reserved.
// Use of this source code is governed by Apache License Version 2.0
// license that can be found in the LICENSE file.

package glean

import (
	"plugin"
	"reflect"
	"runtime"

	"github.com/smallnest/glean/log"
)

// LoadSymbol loads a plugin and gets the symbol.
// It encapsulates plugin.Open and plugin.Lookup methods to a convenient function.
// so is file path of the plugin and name is the symbol.
func LoadSymbol(so, name string) (interface{}, error) {
	p, err := plugin.Open(so)
	if err != nil {
		log.Errorf("failed to open %s: %v", so, err)
		return nil, err
	}
	v, err := p.Lookup(name)
	if err != nil {
		log.Errorf("failed to lookup %s: %v", name, err)
		return nil, err
	}

	return v, nil
}

// Reload loads a function or a variable from the plugin and replace passed function or variable.
// If fails to load, the original function or variable won't be replaced.
func Reload(so, name string, vPtr interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	// a Symbol is a pointer to a variable or function.
	s, err := LoadSymbol(so, name)
	if err != nil {
		return err
	}

	vPtrV := reflect.ValueOf(vPtr)
	if vPtrV.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}

	v := vPtrV.Elem()

	if !v.CanSet() {
		return ErrValueCanNotSet
	}

	v.Set(reflect.ValueOf(s).Elem())
	return nil
}

// ReloadFromPlugin is like Reload but it loads a function or a variable from the given *plugin.Plugin.
func ReloadFromPlugin(p *plugin.Plugin, name string, vPtr interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	s, err := p.Lookup(name)
	if err != nil {
		return err
	}

	vPtrV := reflect.ValueOf(vPtr)
	if vPtrV.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}

	v := vPtrV.Elem()
	if !v.CanSet() {
		return ErrValueCanNotSet
	}

	v.Set(reflect.ValueOf(s).Elem())
	return nil
}
