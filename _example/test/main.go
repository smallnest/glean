package main

import (
	"fmt"
	"time"

	"github.com/smallnest/glean"
	"github.com/smallnest/logi"
)

type AddFunc func(x, y int) int

// test loading functions
func testLoad() {
	var fn AddFunc

	err := glean.Reload("plugins/plugin1/plugin1.so", "Add", &fn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin1: 1+2 = %d\n", fn(1, 2))
	}

	err = glean.Reload("plugins/plugin2/plugin2.so", "Add", &fn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin2: 1+2 = %d\n", fn(1, 2))
	}

	// or

	f, err := glean.LoadSymbol("plugins/plugin1/p1.so", "Add")
	if err != nil {
		fmt.Println(err)
	} else {
		fn = f.(AddFunc)
		fmt.Printf("using plugin1: 1+2 = %d\n", fn(1, 2))
	}

	f, err = glean.LoadSymbol("plugins/plugin2/p2.so", "Add")
	if err != nil {
		fmt.Println(err)
	} else {
		fn = f.(AddFunc)
		fmt.Printf("using plugin2: 1+2 = %d\n", fn(1, 2))
	}

}

// load plugins that has the same location.
func testReplacePlugin() {
	var fn AddFunc

	err := glean.Reload("plugins/plugin1/plugin1.so", "Add", &fn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin1: 1+2 = %d\n", fn(1, 2))
	}

	time.Sleep(time.Minute)

	err = glean.Reload("plugins/plugin1/plugin1.so", "Add", &fn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin1: 1+2 = %d\n", fn(1, 2))
	}

}

func main() {
	var v int

	logi.SetLogger(&logi.DummyLogger{})

	err := glean.Reload("plugins/plugin1/plugin1.so", "V", &v)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin1: value = %d\n", v)
	}

	err = glean.Reload("plugins/plugin2/plugin2.so", "V", &v)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin2: value = %d\n", v)
	}

	// or

	vv, err := glean.LoadSymbol("plugins/plugin1/plugin1.so", "V")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin1: value = %v\n", *(vv.(*int)))
	}

	vv, err = glean.LoadSymbol("plugins/plugin2/plugin2.so", "V")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("using plugin2: value = %v\n", *(vv.(*int)))
	}

}
