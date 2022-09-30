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

func SafeMerge(dst, src interface{}, opts ...Option) (_ interface{}, err error) {
	defer func() {
		if v := recover(); v != nil {
			var ok bool
			if err, ok = v.(error); !ok {
				err = fmt.Errorf("panic: %v", v)
			}
		}
	}()

	return Merge(dst, src, opts...)
}

type merger struct {
	*Options
}

func newMerger(options *Options) *merger {
	return &merger{
		Options: options,
	}
}

func (m *merger) merge(dst, src reflect.Value, resolver Resolver) (
	reflect.Value, error) {
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
	case reflect.Array:
		ret, err = m.mergeArray(dst, src)
	case reflect.Chan:
		ret, err = m.mergeChan(dst, src)
	default:
		// Including reflect.Invalid
		ret, err = m.mergeDefault(dst, src)
	}

	if err != nil {
		return reflect.Value{}, err
	}
	return makeDeepPointer(ret, depth), nil
}

func (m *merger) mergeDefault(dst, src reflect.Value) (reflect.Value, error) {
	return makeValue(src), nil
}

func (m *merger) mergeStruct(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.structStrategy {
	case StructStrategyIgnore:
		ret = makeValue(dst)

	case StructStrategyReplace:
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

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidStrategy, m.structStrategy)
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

	case SliceStrategyRefer:
		ret = makeValue(src)

	case SliceStrategyReplace:
		ret = makeZeroValue(src)
		for i := 0; i < src.Len(); i++ {
			ret = reflect.Append(ret, src.Index(i))
		}

	case SliceStrategyReplaceElementsDynamic:
		max := int(math.Max(float64(dst.Len()), float64(src.Len())))
		ret = makeZeroValue(dst)
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

			dstElem, srcElem, depth, err =
				resolve(dstElem, srcElem, m.sliceResolver)
			if err != nil {
				return reflect.Value{}, err
			}

			if m.conditions.Check(dstElem, srcElem) {
				ret = reflect.Append(ret, makeDeepPointer(srcElem, depth))
			} else {
				ret = reflect.Append(ret, makeDeepPointer(dstElem, depth))
			}
		}

	case SliceStrategyReplaceElementsStatic:
		ret = makeZeroValue(dst)
		for index := 0; index < dst.Len(); index++ {
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

			dstElem, srcElem, depth, err =
				resolve(dstElem, srcElem, m.sliceResolver)
			if err != nil {
				return reflect.Value{}, err
			}

			if m.conditions.Check(dstElem, srcElem) {
				ret = reflect.Append(ret, makeDeepPointer(srcElem, depth))
			} else {
				ret = reflect.Append(ret, makeDeepPointer(dstElem, depth))
			}
		}

	case SliceStrategyReplaceDeepDynamic:
		max := int(math.Max(float64(dst.Len()), float64(src.Len())))
		ret = makeZeroValue(dst)
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
			ret = reflect.Append(ret, v)
		}

	case SliceStrategyReplaceDeepStatic:
		ret = makeZeroValue(dst)
		for index := 0; index < dst.Len(); index++ {
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
			ret = reflect.Append(ret, v)
		}

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidStrategy, m.sliceStrategy)
	}

	return ret, nil
}

func (m *merger) mergeArray(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.arrayStrategy {
	case ArrayStrategyIgnore:
		ret = makeValue(dst)

	case ArrayStrategyReplace:
		ret = makeValue(src)

	case ArrayStrategyReplaceElements:
		ret = makeZeroValue(dst)
		for index := 0; index < dst.Len(); index++ {
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

			dstElem, srcElem, depth, err =
				resolve(dstElem, srcElem, m.arrayResolver)
			if err != nil {
				return reflect.Value{}, err
			}

			if m.conditions.Check(dstElem, srcElem) {
				ret.Index(index).Set(makeDeepPointer(srcElem, depth))
			} else {
				ret.Index(index).Set(makeDeepPointer(dstElem, depth))
			}
		}

	case ArrayStrategyReplaceDeep:
		ret = makeZeroValue(dst)
		for index := 0; index < dst.Len(); index++ {
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
			v, err := m.merge(dstElem, srcElem, m.arrayResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			ret.Index(index).Set(v)
		}

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidStrategy, m.arrayStrategy)
	}

	return ret, nil
}

func (m *merger) mergeChan(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.chanStrategy {
	case ChanStrategyIgnore:
		ret = makeValue(dst)

	case ChanStrategyRefer:
		ret = makeValue(src)

	case ChanStrategyAppend:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)
		retSlice := reflect.AppendSlice(dstSlice, srcSlice)
		ret = sliceToChan(retSlice)

	case ChanStrategyReplace:
		srcSlice := chanToSlice(src)
		ret = sliceToChan(srcSlice)

	case ChanStrategyReplaceElements:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)

		var retSlice = makeZeroValue(dstSlice)
		for index := 0; index < dstSlice.Len(); index++ {
			var (
				dstElem, srcElem reflect.Value
				err              error
				depth            int
			)
			if index < dstSlice.Len() {
				dstElem = dstSlice.Index(index)
			}
			if index < srcSlice.Len() {
				srcElem = srcSlice.Index(index)
			}

			dstElem, srcElem, depth, err =
				resolve(dstElem, srcElem, m.chanResolver)
			if err != nil {
				return reflect.Value{}, err
			}

			if m.conditions.Check(dstElem, srcElem) {
				retSlice = reflect.Append(retSlice,
					makeDeepPointer(srcElem, depth))
			} else {
				retSlice = reflect.Append(retSlice,
					makeDeepPointer(dstElem, depth))
			}
		}

		ret = sliceToChan(retSlice)

	case ChanStrategyReplaceDeep:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)

		retSlice := makeZeroValue(dstSlice)
		for index := 0; index < dstSlice.Len(); index++ {
			var (
				dstElem, srcElem reflect.Value
				err              error
			)
			if index < dstSlice.Len() {
				dstElem = dstSlice.Index(index)
			}
			if index < srcSlice.Len() {
				srcElem = srcSlice.Index(index)
			}

			v, err := m.merge(dstElem, srcElem, m.chanResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			retSlice = reflect.Append(retSlice, v)
		}

		ret = sliceToChan(retSlice)

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidStrategy, m.chanStrategy)
	}

	return ret, nil
}
