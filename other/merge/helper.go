package merge

import (
	"math"
	"reflect"
	"unsafe"
)

func isExported(v reflect.Value, index int) bool {
	return v.Type().Field(index).PkgPath == ""
}

func makeDeepPointer(v reflect.Value, depth int) reflect.Value {
	for index := 0; index < depth; index++ {
		v = makePointer(v)
	}

	return v
}

func makePointer(v reflect.Value) reflect.Value {
	vRet := reflect.New(v.Type())
	vRet.Elem().Set(v)
	return vRet
}

func makeCopiedValue(v reflect.Value) reflect.Value {
	ret := reflect.New(v.Type()).Elem()
	ret.Set(v)
	return ret
}

func makeZeroValue(v reflect.Value) reflect.Value {
	return reflect.New(v.Type()).Elem()
}

func makeEmptyValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(v.Type(), 0, 0)
	case reflect.Map:
		return reflect.MakeMap(v.Type())
	case reflect.Chan:
		return reflect.MakeChan(v.Type(), 0)
	default:
		return reflect.New(v.Type()).Elem()
	}
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
	vSlice := reflect.MakeSlice(
		reflect.SliceOf(vChan.Type().Elem()), 0, vChan.Len())
	for index, length := 0, vChan.Len(); index < length; index++ {
		v, ok := vChan.Recv()
		if !ok {
			break
		}
		vSlice = reflect.Append(vSlice, v)
	}
	for index, length := 0, vSlice.Len(); index < length; index++ {
		vChan.Send(vSlice.Index(index))
	}
	return vSlice
}

func sliceToChan(vSlice reflect.Value) reflect.Value {
	vChan := reflect.MakeChan(
		reflect.ChanOf(reflect.BothDir, vSlice.Type().Elem()), vSlice.Len())
	for index, length := 0, vSlice.Len(); index < length; index++ {
		vChan.Send(vSlice.Index(index))
	}
	return vChan
}

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 |
		uint64 | float32 | float64
}

func max[T number](s ...T) T {
	switch length := len(s); length {
	case 0:
		return 0
	case 1:
		return s[0]
	default:
		max := s[0]
		for index := 1; index < length; index++ {
			max = T(math.Max(float64(max), float64(s[index])))
		}
		return max
	}
}

func min[T number](s ...T) T {
	switch length := len(s); length {
	case 0:
		return 0
	case 1:
		return s[0]
	default:
		max := s[0]
		for index := 1; index < length; index++ {
			max = T(math.Min(float64(max), float64(s[index])))
		}
		return max
	}
}

func maxLen(v ...reflect.Value) int {
	var lengthSlice = make([]int, len(v))
	for index, value := range v {
		lengthSlice[index] = value.Len()
	}
	return max(lengthSlice...)
}

func minLen(v ...reflect.Value) int {
	var lengthSlice = make([]int, len(v))
	for index, value := range v {
		lengthSlice[index] = value.Len()
	}
	return min(lengthSlice...)
}
