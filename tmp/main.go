package main

import (
	"fmt"
	"reflect"
)

func main() {
	type a struct {
		b int
	}
	var c = a{b: 1}
	vc := reflect.ValueOf(c)
	for i := 0; i < vc.NumField(); i++ {
		f := vc.Field(i)
		fmt.Println(f.Int())
	}
}
