package merge

import (
	"fmt"
	"reflect"
)

type Resolver int

const (
	ResolverNone Resolver = iota
	ResolverSingle
	ResolverBoth
	ResolverDeepSingle
	ResolverDeepBoth
)

func resolve(dst, src reflect.Value, resolver Resolver) (
	reflect.Value, reflect.Value, int, error) {
	var depth int
	switch resolver {
	case ResolverNone:
		return dst, src, depth, nil
	case ResolverBoth:
		if dst.Kind() == reflect.Ptr && src.Kind() == reflect.Ptr &&
			dst.Elem().IsValid() && src.Elem().IsValid() {
			dst = dst.Elem()
			src = src.Elem()
			depth++
		}
		return dst, src, depth, nil

	case ResolverDeepBoth:
		for dst.Kind() == reflect.Ptr && src.Kind() == reflect.Ptr &&
			dst.Elem().IsValid() && src.Elem().IsValid() {
			dst = dst.Elem()
			src = src.Elem()
			depth++
		}
		return dst, src, depth, nil

	case ResolverSingle:
		if dst.Kind() == reflect.Ptr && dst.Elem().IsValid() {
			dst = dst.Elem()
			depth++
		}
		if src.Kind() == reflect.Ptr && src.Elem().IsValid() {
			src = src.Elem()
		}
		return dst, src, depth, nil

	case ResolverDeepSingle:
		for dst.Kind() == reflect.Ptr && dst.Elem().IsValid() {
			dst = dst.Elem()
			depth++
		}
		for src.Kind() == reflect.Ptr && src.Elem().IsValid() {
			src = src.Elem()
		}
		return dst, src, depth, nil

	default:
		return reflect.Value{}, reflect.Value{}, depth,
			fmt.Errorf("%w: %v", ErrUnknownResolver, resolver)
	}
}

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
