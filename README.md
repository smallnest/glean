# Glean

*A go plugin framework that can reload variables and functions from plugins automatically.*


[![License](https://img.shields.io/:license-apache-2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/smallnest/glean?status.png)](http://godoc.org/github.com/smallnest/glean)  [![travis](https://travis-ci.org/smallnest/glean.svg?branch=master)](https://travis-ci.org/smallnest/glean) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/glean)](https://goreportcard.com/report/github.com/smallnest/glean) 


## Installation

```sh
go get -u github.com/smallnest/glean
```

## Feature

- load symbol and you don't worry about errors
- load/reload exported variables and funtions from plugins
- watch plugins' changes and reload pointer of variables and function in applications

**Notice** glean only can reload functions or variables that can be addresses.

## Examples

see [Examples](https://github.com/smallnest/glean/tree/master/_example)

Let's look the httpserver example to learn how to use glean.

### httpserver

httpserver is a very very simple http server.

A simple http server is just like this:

```go
var FooHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}

http.Handle("/foo", fooHandler)

log.Fatal(http.ListenAndServe(":9988", nil))
```

Our goal is to replace fooHandler with latest code dynamically (hot fix).

```go
var FooHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, gp")
}
```

No need to restart this server.

**step1:** build the two plugin

enter `_example/httpserver/plugins/plugin1` and `_example/httpserver/plugins/plugin2`, and run the `build.sh` to generate the so file.

Currently plugin supports linux and MacOS.

**step2:** modify the server implementation

```go
package main

import (
	"log"
	"net/http"

	"github.com/smallnest/glean"
)

func main() {
	g := glean.New("plugin.json")
	err := g.LoadConfig()
	if err != nil {
		panic(err)
	}

	var fooHandler func(w http.ResponseWriter, r *http.Request)

	err = g.ReloadAndWatch("FooHandlerID", &fooHandler)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/foo", WarpFuncPtr(&fooHandler))

	log.Fatal(http.ListenAndServe(":9988", nil))
}

func WarpFuncPtr(fn *func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		(*fn)(w, r)
	}
}
```

Firstly create the Glean instance and load config from given file.
And then use `ReloadAndWatch` to load fooHandler and begin to watch its changes.
At last use WarpFuncPtr to wrap fooHandler as a HandleFunc.

Run `go run main.go` to start this server, use a browser to visit "http://locakhost:9988/foo" and you will see `hello world`

Change the config file `plugin.json` and replace `"file": "plugins/plugin1/plugin1.so"` with :

```
"file": "plugins/plugin2/plugin2.so",
```

Browser the prior location and you will see `hello gp`.
