package merge

import (
	"errors"
	"reflect"
)

var (
	ErrKindNotSupported = errors.New("must be a struct, map or slice")
	ErrNotAdrressable   = errors.New("must not be unadrressable value")
	ErrNilValue         = errors.New("must not be nil value")
	ErrInvalidValue     = errors.New("must not be invalid value")
	ErrUnknownRange     = errors.New("unknown range")
)

func resolveAddressableContainer(v reflect.Value) (reflect.Value, error) {
	if !v.CanAddr() {
		return reflect.Value{}, ErrNotAdrressable
	}

	return resolveContainer(v)
}

func resolveContainer(v reflect.Value) (reflect.Value, error) {
	var err error

	v, err = resolve(v)
	if err != nil {
		return reflect.Value{}, err
	}

	if v.Kind() != reflect.Struct && v.Kind() != reflect.Map && v.Kind() != reflect.Slice {
		return reflect.Value{}, ErrKindNotSupported
	}

	return v, nil
}

func resolveAdrressable(v reflect.Value) (reflect.Value, error) {
	if !v.CanSet() {
		return reflect.Value{}, ErrNotAdrressable
	}

	return resolve(v)
}

func resolve(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return v, ErrInvalidValue
	}
	for {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return v, ErrNilValue
			}
			v = v.Elem()
		} else {
			break
		}
	}
	return v, nil
}
