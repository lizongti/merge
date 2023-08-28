package merge

import (
	"fmt"
	"reflect"
	"unicode"
)

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMap(dst, src reflect.Value, visited map[uintptr]*visit, depth int, config *Options) (err error) {
	overwrite := config.overwrite
	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return nil
			}
		}
		// Remember, remember...
		visited[h] = &visit{typ, seen, addr}
	}
	zeroValue := reflect.Value{}
	switch dst.Kind() {
	case reflect.Map:
		dstMap := dst.Interface().(map[string]interface{})
		for i, n := 0, src.NumField(); i < n; i++ {
			srcType := src.Type()
			field := srcType.Field(i)
			if !isExported(field) {
				continue
			}
			fieldName := field.Name
			fieldName = changeInitialCase(fieldName, unicode.ToLower)
			if v, ok := dstMap[fieldName]; !ok || (isEmptyValue(reflect.ValueOf(v)) || overwrite) {
				dstMap[fieldName] = src.Field(i).Interface()
			}
		}
	case reflect.Ptr:
		if dst.IsNil() {
			v := reflect.New(dst.Type().Elem())
			dst.Set(v)
		}
		dst = dst.Elem()
		fallthrough
	case reflect.Struct:
		srcMap := src.Interface().(map[string]interface{})
		for key := range srcMap {
			config.overwriteWithEmptyValue = true
			srcValue := srcMap[key]
			fieldName := changeInitialCase(key, unicode.ToUpper)
			dstElement := dst.FieldByName(fieldName)
			if dstElement == zeroValue {
				// We discard it because the field doesn't exist.
				continue
			}
			srcElement := reflect.ValueOf(srcValue)
			dstKind := dstElement.Kind()
			srcKind := srcElement.Kind()
			if srcKind == reflect.Ptr && dstKind != reflect.Ptr {
				srcElement = srcElement.Elem()
				srcKind = reflect.TypeOf(srcElement.Interface()).Kind()
			} else if dstKind == reflect.Ptr {
				// Can this work? I guess it can't.
				if srcKind != reflect.Ptr && srcElement.CanAddr() {
					srcPtr := srcElement.Addr()
					srcElement = reflect.ValueOf(srcPtr)
					srcKind = reflect.Ptr
				}
			}

			if !srcElement.IsValid() {
				continue
			}
			if srcKind == dstKind {
				if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
					return
				}
			} else if dstKind == reflect.Interface && dstElement.Kind() == reflect.Interface {
				if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
					return
				}
			} else if srcKind == reflect.Map {
				if err = deepMap(dstElement, srcElement, visited, depth+1, config); err != nil {
					return
				}
			} else {
				return fmt.Errorf("type mismatch on %s field: found %v, expected %v", fieldName, srcKind, dstKind)
			}
		}
	}
	return
}

// Map sets fields' values in dst from src.
// src can be a map with string keys or a struct. dst must be the opposite:
// if src is a map, dst must be a valid pointer to struct. If src is a struct,
// dst must be map[string]interface{}.
// It won't merge unexported (private) fields and will do recursively
// any exported field.
// If dst is a map, keys will be src fields' names in lower camel case.
// Missing key in src that doesn't match a field in dst will be skipped. This
// doesn't apply if dst is a map.
// This is separated method from Merge because it is cleaner and it keeps sane
// semantics: merging equal types, mapping different (restricted) types.
func Map(dst, src interface{}, opts ...Option) error {
	return _map(dst, src, opts...)
}

func _map(dst, src interface{}, opts ...Option) error {
	if dst != nil && reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return ErrNonPointerAgument
	}
	var (
		vDst, vSrc reflect.Value
		err        error
	)
	config := &Options{}

	for _, opt := range opts {
		opt(config)
	}

	if vDst, vSrc, err = resolveValues(dst, src); err != nil {
		return err
	}
	// To be friction-less, we redirect equal-type arguments
	// to deepMerge. Only because arguments can be anything.
	if vSrc.Kind() == vDst.Kind() {
		return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, config)
	}
	switch vSrc.Kind() {
	case reflect.Struct:
		if vDst.Kind() != reflect.Map {
			return ErrExpectedMapAsDestination
		}
	case reflect.Map:
		if vDst.Kind() != reflect.Struct {
			return ErrExpectedStructAsDestination
		}
	default:
		return ErrNotSupported
	}
	return deepMap(vDst, vSrc, make(map[uintptr]*visit), 0, config)
}
