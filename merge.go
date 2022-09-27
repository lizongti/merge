package merge

import (
	"fmt"
	"reflect"
)

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...Option) (interface{}, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)

	vRet, err := newMerger(options).merge(vDst, vSrc)
	if err != nil {
		return nil, err
	}

	return vRet.Interface(), nil
}

func MustMerge(dst, src interface{}, opts ...Option) interface{} {
	v, err := Merge(dst, src, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

type merger struct {
	*Options
}

func newMerger(options *Options) *merger {
	return &merger{
		Options: options,
	}
}

func (m *merger) merge(dst, src reflect.Value) (reflect.Value, error) {
	dst, src, depth, err := resolve(dst, src, m.resolver)
	if err != nil {
		return reflect.Value{}, err
	}

	var vRet reflect.Value
	kindGroup := getKindGroup(dst)
	switch kindGroup {
	case kindGroupContainer:
		vRet, err = m.mergeContainer(dst, src)
		if err != nil {
			return reflect.Value{}, err
		}
		return makePointerInDepth(vRet, depth), nil
	case kindGroupRefer:
		vRet, err = m.mergeRefer(dst, src)
		if err != nil {
			return reflect.Value{}, err
		}
		return makePointerInDepth(vRet, depth), nil
	case kindGroupValue:
		vRet, err = m.mergeValue(dst, src)
		if err != nil {
			return reflect.Value{}, err
		}
		return makePointerInDepth(vRet, depth), nil
	default:
		return reflect.Value{}, ErrInvalidValue
	}
}

func (m *merger) mergeContainer(dst, src reflect.Value) (reflect.Value, error) {
	switch dst.Kind() {
	case reflect.Struct:
		return m.mergeStruct(dst, src)
	case reflect.Map:
		return m.mergeMap(dst, src)
	case reflect.Slice:
		return m.mergeSlice(dst, src)
	case reflect.Array:
		panic("not implemented")
	case reflect.Chan:
		panic("not implemented")
	default:
		return reflect.Value{},
			fmt.Errorf("%w: %s", ErrKindNotSupported, dst.Kind())
	}
}

func (m *merger) mergeRefer(dst, src reflect.Value) (reflect.Value, error) {
	if !m.conditions.canCover(dst, src) {
		return dst, nil
	}

	vRet := reflect.New(src.Type()).Elem()
	vRet.Set(src)
	return vRet, nil
}

func (m *merger) mergeValue(dst, src reflect.Value) (reflect.Value, error) {
	if !m.conditions.canCover(dst, src) {
		return dst, nil
	}

	vRet := reflect.New(src.Type()).Elem()
	vRet.Set(src)
	return vRet, nil
}

func (m *merger) mergeStruct(dst, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func (m *merger) mergeMap(dst, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func (m *merger) mergeSlice(dst, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}
