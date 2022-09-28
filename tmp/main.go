package main

import (
	"fmt"
	"reflect"

	"github.com/cloudlibraries/merge/t"
)

type B struct {
	C int
}
type A struct {
	B
	d int
}

func main() {
	var f t.F
	// v := reflect.TypeOf(a)
	// for i := 0; i < v.NumField(); i++ {
	// 	field := v.Field(i)
	// 	fmt.Println(field.Name, field.Anonymous, ast.IsExported(field.Name), len(field.PkgPath) == 0)
	// }
	fmt.Println(getUnexportedField(reflect.ValueOf(f), "g"))
}

func getUnexportedField(v reflect.Value, name string) reflect.Value {
	return v.FieldByNameFunc(func(s string) bool {
		return s == name
	})
}
