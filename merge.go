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

func (m *merger) merge(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	var err error

	if dst, err = m.resolve(dst); err != nil {
		return reflect.Value{}, err
	}

	if src, err = m.resolve(src); err != nil {
		return reflect.Value{}, err
	}
	dst.IsZero()

	// if dst.Kind() != src.Kind() && m.mappingEnabled {
	// 	return m.mapping(dst, src)
	// }

	return m.deepMerge(dst, src)
}

func (m *merger) resolve(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return v, ErrInvalidValue
	}
	switch m.resolver {
	case ResolverNone:
		return v, nil
	case ResolverNormal:
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return v, ErrNilValue
			}
			v = v.Elem()
		}
		return v, nil
	case ResolverDeep:
		for {
			if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
				if v.IsNil() {
					return v, ErrNilValue
				}
				v = v.Elem()
			} else {
				break
			}
		}
		return v, nil
	default:
		return v, fmt.Errorf("%w: %v", ErrUnknownResolver, m.resolver)
	}
}

func (m *merger) deepMerge(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	kindGroup := resolveKindGroup(dst)
	switch kindGroup {
	case kindGroupContainer:
		switch dst.Kind() {
		case reflect.Struct:
			return m.deepMergeStruct(dst, src)
		case reflect.Map:
			return m.deepMergeMap(dst, src)
		case reflect.Slice:
			return m.deepMergeSlice(dst, src)
		}
	case kindGroupReference:
		return m.deepMergeReference(dst, src)
	case kindGroupValue:
		return m.deepMergeValue(dst, src)
	default:
		return reflect.Value{}, ErrInvalidValue
	}
}

func (m *merger) deepMergeStruct(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func (m *merger) deepMergeMap(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func (m *merger) deepMergeSlice(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func (m *merger) deepMergeReference(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	// resolve
	// vRet := reflect.New(src.Type()).Elem()
	// vRet.Set(src)
	return reflect.Value{}, nil
}

func (m *merger) deepMergeValue(dst reflect.Value, src reflect.Value) (reflect.Value, error) {
	vRet := reflect.New(src.Type()).Elem()
	vRet.Set(src)
	return vRet, nil
}
