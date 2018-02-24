// Copyright 2009 smallnest. All rights reserved.
// Use of this source code is governed by Apache License Version 2.0
// license that can be found in the LICENSE file.

package glean

import (
	"plugin"
	"testing"

	"github.com/smallnest/glean/log"
)

func TestLoadSymbol(t *testing.T) {
	log.SetDummyLogger()

	type args struct {
		so   string
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"Add",
			args{
				"_example/test/plugins/plugin1/plugin1.so",
				"Add",
			},
			nil,
			false,
		},
		{
			"v",
			args{
				"_example/test/plugins/plugin1/plugin1.so",
				"v",
			},
			nil,
			true,
		},
		{
			"nonExisted",
			args{
				"_example/test/plugins/pluginabc/plugin1.so",
				"Add",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadSymbol(tt.args.so, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("LoadSymbol() want nil but got: %v", got)
			}
		})
	}
}

func TestReload(t *testing.T) {
	var fn func(x, y int) int
	var v int

	type args struct {
		so   string
		name string
		vPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Add",
			args{
				"_example/test/plugins/plugin1/plugin1.so",
				"Add",
				&fn,
			},
			false,
		},
		{
			"V",
			args{
				"_example/test/plugins/plugin1/plugin1.so",
				"V",
				&v,
			},
			false,
		},
		{
			"v",
			args{
				"_example/test/plugins/plugin1/plugin1.so",
				"v",
				&v,
			},
			true,
		},
		{
			"nonExisted",
			args{
				"_example/test/plugins/pluginabc/plugin1.so",
				"Add",
				&fn,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Reload(tt.args.so, tt.args.name, tt.args.vPtr); (err != nil) != tt.wantErr {
				t.Errorf("Reload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReloadFromPlugin(t *testing.T) {
	var fn func(x, y int) int
	var v int

	p, err := plugin.Open("_example/test/plugins/plugin1/plugin1.so")
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	type args struct {
		p    *plugin.Plugin
		name string
		vPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Add",
			args{
				p,
				"Add",
				&fn,
			},
			false,
		},
		{
			"V",
			args{
				p,
				"V",
				&v,
			},
			false,
		},
		{
			"v",
			args{
				p,
				"v",
				&v,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReloadFromPlugin(tt.args.p, tt.args.name, tt.args.vPtr); (err != nil) != tt.wantErr {
				t.Errorf("ReloadFromPlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
