// Copyright 2009 smallnest. All rights reserved.
// Use of this source code is governed by Apache License Version 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
)

var FooHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
