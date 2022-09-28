package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type b struct {
	c int
}
type A struct {
	b
	d int
}

func main() {
	// GetUnexportedField(reflect.ValueOf(&Foo{}).Elem().FieldByName("unexportedField"))
	var f A
	// ret := reflect.New(reflect.ValueOf(f).Type()).Elem()
	ret := reflect.ValueOf(f)
	field := ret.FieldByName("b")
	// v := reflect.TypeOf(a)
	// for i := 0; i < v.NumField(); i++ {
	// 	field := v.Field(i)
	// 	fmt.Println(field.Name, field.Anonymous, ast.IsExported(field.Name), len(field.PkgPath) == 0)
	// }
	// fmt.Println(getUnexportedField(reflect.ValueOf(f), "g"))

	SetUnexportedField(field, b{c: 1})
	// *(reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Interface().(*int)) = 1

	fmt.Println(f)
}

type Foo struct {
	unexportedField string
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}
