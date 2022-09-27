package main

import (
	"fmt"
	"reflect"
)

func main() {
	// var a []int = []int{1, 2, 3}
	var b = reflect.Value{}
	// va := reflect.ValueOf(a)
	fmt.Println(b.Kind())

}
