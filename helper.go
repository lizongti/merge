package merge

import (
	"errors"
	"reflect"
	"unsafe"
)

var (
	ErrKindNotSupported = errors.New("kind not supported")
	ErrNotAdrressable   = errors.New("must not be unadrressable value")
	ErrNilValue         = errors.New("must not be nil value")
	ErrInvalidValue     = errors.New("must not be invalid value")
	ErrUnknownRange     = errors.New("unknown range")
	ErrUnknownResolver  = errors.New("unknown resolver")
	ErrNotSettable      = errors.New("must be settable")
	ErrInvalidStrategy  = errors.New("invalid strategy")
)

func makeDeepPointer(v reflect.Value, depth int) reflect.Value {
	for i := 0; i < depth; i++ {
		v = makePointer(v)
	}

	return v
}

func makePointer(v reflect.Value) reflect.Value {
	vRet := reflect.New(v.Type())
	vRet.Elem().Set(v)
	return vRet
}

func makeValue(v reflect.Value) reflect.Value {
	ret := reflect.New(v.Type()).Elem()
	ret.Set(v)
	return ret
}

func makeZeroValue(v reflect.Value) reflect.Value {
	return reflect.New(v.Type()).Elem()
}

func getFieldByName(v reflect.Value, name string) reflect.Value {
	return v.FieldByNameFunc(func(s string) bool {
		return s == name
	})
}

func setValueToField(field reflect.Value, value reflect.Value) {
	unsafePtr := unsafe.Pointer(field.UnsafeAddr())
	reflect.NewAt(field.Type(), unsafePtr).Elem().Set(value)
}

func getValueFromField(field reflect.Value) reflect.Value {
	unsafePtr := unsafe.Pointer(field.UnsafeAddr())
	return reflect.NewAt(field.Type(), unsafePtr).Elem()
}

func chanToSlice(vChan reflect.Value) reflect.Value {
	vSlice := reflect.MakeSlice(vChan.Type().Elem(), 0, vChan.Len())
	for i, n := 0, vChan.Len(); i < n; i++ {
		v, ok := vChan.Recv()
		if !ok {
			break
		}
		vSlice.Set(reflect.Append(vSlice, v))
	}
	for i := 0; i < vSlice.Len(); i++ {
		vChan.Send(vSlice.Index(i))
	}
	return vSlice
}

func sliceToChan(vSlice reflect.Value) reflect.Value {
	vChan := reflect.MakeChan(vSlice.Type().Elem(), vSlice.Len())
	for i := 0; i < vSlice.Len(); i++ {
		vChan.Send(vSlice.Index(i))
	}
	return vChan
}
