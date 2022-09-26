package main

import (
	"fmt"
	"reflect"

	"github.com/cloudlibraries/merge"
)

func set(a interface{}) {
	v := reflect.ValueOf(a)
	b := 20
	v.Set(reflect.ValueOf(&b))
}

func main() {

	fmt.Println(merge.Merge(10, 1))
}
