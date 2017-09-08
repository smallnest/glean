// Copyright 2009 smallnest. All rights reserved.
// Use of this source code is governed by Apache License Version 2.0
// license that can be found in the LICENSE file.

package glean

import "testing"

func TestGP_LoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		g       *Glean
		wantErr bool
	}{
		{
			name:    "normal",
			g:       New("./example/plugin.json"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.LoadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Glean.LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			var fn func(x, y int) int
			var v int

			err := tt.g.Reload("EF5A35EC-46EB-4E62-8251-78F1A49FA7DC", &fn)
			if err != nil {
				t.Errorf("failed to reload fn: %v", err)
			}

			err = tt.g.Reload("2E8FD057-99EC-41B9-8172-0EBF18F9A48D", &v)
			if err != nil {
				t.Errorf("failed to reload v: %v", err)
			}
		})
	}
}
