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
	case reflect.Array:
		ret, err = m.mergeArray(dst, src)
	case reflect.Struct:
		ret, err = m.mergeStruct(dst, src)
	case reflect.Slice:
		ret, err = m.mergeSlice(dst, src)
	case reflect.Chan:
		ret, err = m.mergeChan(dst, src)
	case reflect.Map:
		ret, err = m.mergeMap(dst, src)
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
	return makeCopiedValue(src), nil
}

func (m *merger) mergeArray(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.arrayStrategy {
	case ArrayStrategyIgnore:
		ret = makeCopiedValue(dst)

	case ArrayStrategyReplace:
		ret = makeCopiedValue(src)

	case ArrayStrategyReplaceElem:
		ret = makeZeroValue(dst)
		for index, length := 0, maxLen(dst); index < length; index++ {
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
		for index, length := 0, maxLen(dst); index < length; index++ {
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

func (m *merger) mergeStruct(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.structStrategy {
	case StructStrategyIgnore:
		ret = makeCopiedValue(dst)

	case StructStrategyReplace:
		ret = makeCopiedValue(src)

	case StructStrategyReplaceElem:
		src = makeCopiedValue(src)
		dst = makeCopiedValue(dst)
		ret = makeZeroValue(dst)
		for index, length := 0, dst.NumField(); index < length; index++ {
			dstValue := getValueFromField(dst.Field(index))
			srcValue := getValueFromField(src.Field(index))
			if m.conditions.Check(dstValue, srcValue) {
				setValueToField(ret.Field(index), srcValue)
			} else {
				setValueToField(ret.Field(index), dstValue)
			}
		}

	case StructStrategyReplaceDeep:
		src = makeCopiedValue(src)
		dst = makeCopiedValue(dst)
		ret = makeZeroValue(dst)
		for index, length := 0, dst.NumField(); index < length; index++ {
			dstValue := getValueFromField(dst.Field(index))
			srcValue := getValueFromField(src.Field(index))
			v, err := m.merge(dstValue, srcValue, m.structResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			setValueToField(ret.Field(index), v)
		}

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidStrategy, m.structStrategy)
	}

	return ret, nil
}

func (m *merger) mergeSlice(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.sliceStrategy {
	case SliceStrategyIgnore:
		ret = makeCopiedValue(dst)

	case SliceStrategyAppend:
		ret = makeCopiedValue(reflect.AppendSlice(dst, src))

	case SliceStrategyRefer:
		ret = makeCopiedValue(src)

	case SliceStrategyReplace:
		ret = makeEmptyValue(src)
		for index, length := 0, maxLen(src); index < length; index++ {
			ret = reflect.Append(ret, src.Index(index))
		}

	case SliceStrategyReplaceElem:
		ret = makeEmptyValue(dst)
		for index, length := 0, maxLen(dst); index < length; index++ {
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

	case SliceStrategyReplaceElemDynamic:
		ret = makeEmptyValue(dst)
		for index, length := 0, maxLen(src, dst); index < length; index++ {
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

	case SliceStrategyReplaceDeep:
		ret = makeEmptyValue(dst)
		for index, length := 0, maxLen(dst); index < length; index++ {
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

	case SliceStrategyReplaceDeepDynamic:
		ret = makeEmptyValue(dst)
		for index, length := 0, maxLen(src, dst); index < length; index++ {
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

func (m *merger) mergeChan(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.chanStrategy {
	case ChanStrategyIgnore:
		ret = makeCopiedValue(dst)

	case ChanStrategyRefer:
		ret = makeCopiedValue(src)

	case ChanStrategyAppend:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)
		retSlice := reflect.AppendSlice(dstSlice, srcSlice)
		ret = sliceToChan(retSlice)

	case ChanStrategyReplace:
		srcSlice := chanToSlice(src)
		ret = sliceToChan(srcSlice)

	case ChanStrategyReplaceElem:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)

		retSlice := makeEmptyValue(dstSlice)
		for index, length := 0, maxLen(dstSlice); index < length; index++ {
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

	case ChanStrategyReplaceElemDynamic:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)

		retSlice := makeEmptyValue(dstSlice)
		for index, length := 0,
			maxLen(srcSlice, dstSlice); index < length; index++ {
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

		retSlice := makeEmptyValue(dstSlice)
		for index, length := 0, maxLen(dstSlice); index < length; index++ {
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

	case ChanStrategyReplaceDeepDynamic:
		dstSlice := chanToSlice(dst)
		srcSlice := chanToSlice(src)

		retSlice := makeEmptyValue(dstSlice)
		for index, length := 0,
			maxLen(srcSlice, dstSlice); index < length; index++ {
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

func (m *merger) mergeMap(dst, src reflect.Value) (reflect.Value, error) {
	var ret reflect.Value

	switch m.mapStrategy {
	case MapStrategyIgnore:
		ret = makeCopiedValue(dst)

	case MapStrategyRefer:
		ret = makeCopiedValue(src)

	case MapStrategyReplace:
		ret = makeEmptyValue(dst)

		for _, key := range src.MapKeys() {
			ret.SetMapIndex(key, src.MapIndex(key))
		}

	case MapStrategyReplaceElem:
		ret = makeEmptyValue(dst)

		for _, key := range dst.MapKeys() {
			srcValue := src.MapIndex(key)
			dstValue := dst.MapIndex(key)
			if m.conditions.Check(dstValue, srcValue) {
				ret.SetMapIndex(key, srcValue)
			} else {
				ret.SetMapIndex(key, dstValue)
			}
		}

	case MapStrategyReplaceElemDynamic:
		ret = makeEmptyValue(dst)

		for _, key := range dst.MapKeys() {
			srcValue := src.MapIndex(key)
			dstValue := dst.MapIndex(key)
			if m.conditions.Check(dstValue, srcValue) {
				ret.SetMapIndex(key, srcValue)
			} else {
				ret.SetMapIndex(key, dstValue)
			}
		}

		for _, key := range src.MapKeys() {
			if ret.MapIndex(key).IsValid() {
				continue
			}

			srcValue := src.MapIndex(key)
			dstValue := dst.MapIndex(key)
			if m.conditions.Check(dstValue, srcValue) {
				ret.SetMapIndex(key, srcValue)
			} else {
				ret.SetMapIndex(key, dstValue)
			}
		}

	case MapStrategyReplaceDeep:
		ret = makeEmptyValue(dst)

		for _, key := range dst.MapKeys() {
			srcValue := src.MapIndex(key)
			dstValue := dst.MapIndex(key)
			v, err := m.merge(dstValue, srcValue, m.mapResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			ret.SetMapIndex(key, v)
		}

	case MapStrategyReplaceDeepDynamic:
		ret = makeEmptyValue(dst)

		for _, key := range dst.MapKeys() {
			srcValue := src.MapIndex(key)
			dstValue := dst.MapIndex(key)
			v, err := m.merge(dstValue, srcValue, m.mapResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			ret.SetMapIndex(key, v)
		}

		for _, key := range src.MapKeys() {
			if ret.MapIndex(key).IsValid() {
				continue
			}

			srcValue := src.MapIndex(key)
			dstValue := dst.MapIndex(key)
			v, err := m.merge(dstValue, srcValue, m.mapResolver)
			if err != nil {
				return reflect.Value{}, err
			}
			ret.SetMapIndex(key, v)
		}

	default:
		return reflect.Value{},
			fmt.Errorf("%w: %v", ErrInvalidStrategy, m.mapStrategy)
	}

	return ret, nil
}
