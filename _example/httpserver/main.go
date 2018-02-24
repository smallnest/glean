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
