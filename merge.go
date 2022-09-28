package merge

import (
	"fmt"
	"math"
	"reflect"
)

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...Option) (interface{}, error) {
	options := newOptions(opts)
	merger := newMerger(options)

	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)

	vRet, err := merger.merge(vDst, vSrc, merger.defaultResolver)
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

func (m *merger) merge(dst, src reflect.Value, resolver Resolver) (reflect.Value, error) {
	dst, src, depth, err := resolve(dst, src, resolver)
	if err != nil {
		return reflect.Value{}, err
	}

	if !m.conditions.Check(dst, src) {
		return dst, nil
	}

	var ret reflect.Value
	switch dst.Kind() {
	case reflect.Map:
		ret, err = m.mergeMap(dst, src)
	case reflect.Slice:
		ret, err = m.mergeSlice(dst, src)
	case reflect.Struct:
		ret, err = m.mergeStruct(dst, src)
	// case reflect.Array:
	// 	vRet, err = m.mergeArray(dst, src)
	// case reflect.Chan:
	// 	vRet, err = m.mergeChan(dst, src)
	default:
		// Including reflect.Invalid
		ret, err = m.mergeValue(dst, src)
	}

	if err != nil {
		return reflect.Value{}, err
	}
	return makePointerInDepth(ret, depth), nil
}

func (m *merger) mergeValue(dst, src reflect.Value) (reflect.Value, error) {
	return makeValue(src), nil
}

func (m *merger) mergeStruct(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.structStrategy {
	case StructStrategyIgnore:
		ret = makeValue(dst)

	case StructStrategyReplaceStruct:
		ret = makeValue(src)

	case StructStrategyReplaceFields:
		src = makeValue(src)
		dst = makeValue(dst)
		ret = makeZeroValue(dst)
		for i := 0; i < dst.NumField(); i++ {
			dstValue := getValueFromField(dst.Field(i))
			srcValue := getValueFromField(src.Field(i))
			if m.conditions.Check(dstValue, srcValue) {
				setValueToField(ret.Field(i), srcValue)
			} else {
				setValueToField(ret.Field(i), dstValue)
			}
		}

	case StructStrategyReplaceDeep:
		src = makeValue(src)
		dst = makeValue(dst)
		ret = makeZeroValue(dst)
		for i := 0; i < dst.NumField(); i++ {
			dstValue := getValueFromField(dst.Field(i))
			srcValue := getValueFromField(src.Field(i))
			v, err := m.merge(dstValue, srcValue, m.structResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			setValueToField(ret.Field(i), v)
		}
	}

	return ret, nil
}

func (m *merger) mergeMap(dst, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func (m *merger) mergeSlice(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.sliceStrategy {
	case SliceStrategyIgnore:
		ret = makeValue(dst)

	case SliceStrategyAppend:
		ret = makeValue(reflect.AppendSlice(dst, src))

	case SliceStrategyReplaceSlice:
		ret = makeValue(src)

	case SliceStrategyReplaceElements:
		ret = makeZeroValue(dst)
		max := int(math.Max(float64(dst.Len()), float64(src.Len())))
		for index := 0; index < max; index++ {
			var (
				dstElem, srcElem reflect.Value
				err              error
				depth            int
			)
			if index < dst.Len() {
				dstElem = dst.Index(index)
			}
			if index < src.Len() {
				srcElem = src.Index(index)
			}

			dstElem, srcElem, depth, err = resolve(dstElem, srcElem, m.defaultResolver)
			if err != nil {
				return reflect.Value{}, err
			}

			if m.conditions.Check(dstElem, srcElem) {
				ret = reflect.Append(ret, makePointerInDepth(srcElem, depth))
			} else {
				ret = reflect.Append(ret, makePointerInDepth(dstElem, depth))
			}
		}

	case SliceStrategyReplaceDeep:
		ret = makeZeroValue(dst)
		max := int(math.Max(float64(dst.Len()), float64(src.Len())))
		for index := 0; index < max; index++ {
			var (
				dstElem, srcElem reflect.Value
				err              error
			)
			if index < dst.Len() {
				dstElem = dst.Index(index)
			}
			if index < src.Len() {
				srcElem = src.Index(index)
			}
			v, err := m.merge(dstElem, srcElem, m.sliceResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			ret.Set(reflect.Append(ret, v))
		}

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidSliceStrategy, m.sliceStrategy)
	}

	return ret, nil
}
