package merge

import (
	"errors"
	"reflect"
)

var (
	ErrKindNotSupported     = errors.New("kind not supported")
	ErrNotAdrressable       = errors.New("must not be unadrressable value")
	ErrNilValue             = errors.New("must not be nil value")
	ErrInvalidValue         = errors.New("must not be invalid value")
	ErrUnknownRange         = errors.New("unknown range")
	ErrUnknownResolver      = errors.New("unknown resolver")
	ErrNotSettable          = errors.New("must be settable")
	ErrInvalidSliceStrategy = errors.New("invalid slice strategy")
)

func makePointerInDepth(v reflect.Value, depth int) reflect.Value {
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
