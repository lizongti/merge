package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a = [2]int{}
	reflect.ValueOf(&a).Elem().Index(0).SetInt(1)
	fmt.Println(a)
}
