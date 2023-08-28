package merge

import (
	"fmt"
	"reflect"
)

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMerge(dst, src reflect.Value, visited map[uintptr]*visit, depth int, config *Options) (err error) {
	transformers := config.transformers

	overwrite := config.overwrite
	appendSlice := config.appendSlice
	overwriteWithEmptyValue := config.overwriteWithEmptyValue
	overwriteSliceWithEmptyValue := config.overwriteSliceWithEmptyValue
	// recursiveOverwrite := config.recursiveOverwrite

	if !src.IsValid() {
		return
	}
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
		visited[h] = &visit{typ, seen, addr}
	}

	if transformers != nil && !isReflectNil(dst) && dst.IsValid() {
		if fn := transformers.Transformer(dst.Type()); fn != nil {
			err = fn(dst, src)
			return
		}
	}

	switch dst.Kind() {
	case reflect.Struct:
		if hasMergeableFields(dst) {
			for i, n := 0, dst.NumField(); i < n; i++ {
				if err = deepMerge(dst.Field(i), src.Field(i), visited, depth+1, config); err != nil {
					return
				}
			}
		} else {
			if dst.CanSet() && (isReflectNil(dst) || overwrite) && (!isEmptyValue(src) || overwriteWithEmptyValue) {
				dst.Set(src)
			}
		}
	case reflect.Map:
		if dst.IsNil() && !src.IsNil() {
			if dst.CanSet() {
				dst.Set(reflect.MakeMap(dst.Type()))
			} else {
				dst = src
				return
			}
		}

		if src.Kind() != reflect.Map {
			if overwrite {
				dst.Set(src)
			}
			return
		}

		for _, key := range src.MapKeys() {
			srcElement := src.MapIndex(key)
			if !srcElement.IsValid() {
				continue
			}
			dstElement := dst.MapIndex(key)
			switch srcElement.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Slice:
				if srcElement.IsNil() {
					if overwriteWithEmptyValue || overwriteSliceWithEmptyValue {
						dst.SetMapIndex(key, srcElement)
					}
					continue
				}
				fallthrough
			default:
				if !srcElement.CanInterface() {
					continue
				}
				switch reflect.TypeOf(srcElement.Interface()).Kind() {
				case reflect.Struct:
					fallthrough
				case reflect.Ptr:
					fallthrough
				case reflect.Map:
					srcMapElm := srcElement
					dstMapElm := dstElement
					if srcMapElm.CanInterface() {
						srcMapElm = reflect.ValueOf(srcMapElm.Interface())
						if dstMapElm.IsValid() {
							dstMapElm = reflect.ValueOf(dstMapElm.Interface())
						}
					}
					if err = deepMerge(dstMapElm, srcMapElm, visited, depth+1, config); err != nil {
						return
					}
				case reflect.Slice:
					srcSlice := reflect.ValueOf(srcElement.Interface())

					var dstSlice reflect.Value
					if !dstElement.IsValid() || dstElement.IsNil() {
						dstSlice = reflect.MakeSlice(srcSlice.Type(), 0, srcSlice.Len())
					} else {
						dstSlice = reflect.ValueOf(dstElement.Interface())
					}

					if (!isEmptyValue(srcSlice) || overwriteWithEmptyValue || overwriteSliceWithEmptyValue) && (overwrite || isEmptyValue(dstSlice)) && !appendSlice {
						dstSlice = srcSlice
					} else if appendSlice {
						if srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot append two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = reflect.AppendSlice(dstSlice, srcSlice)
					}
					dst.SetMapIndex(key, dstSlice)
				}
			}
			if dstElement.IsValid() && !isEmptyValue(dstElement) && (reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Map || reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Slice) {
				continue
			}

			if srcElement.IsValid() && ((srcElement.Kind() != reflect.Ptr && overwrite) || !dstElement.IsValid() || isEmptyValue(dstElement)) {
				if dst.IsNil() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				dst.SetMapIndex(key, srcElement)
			}
		}
	case reflect.Slice:
		if !dst.CanSet() {
			break
		}
		if (!isEmptyValue(src) || overwriteWithEmptyValue || overwriteSliceWithEmptyValue) && (overwrite || isEmptyValue(dst)) && !appendSlice {
			dst.Set(src)
		} else if appendSlice {
			if src.Type() != dst.Type() {
				return fmt.Errorf("cannot append two slice with different type (%s, %s)", src.Type(), dst.Type())
			}
			dst.Set(reflect.AppendSlice(dst, src))
		}
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if isReflectNil(src) {
			if overwriteWithEmptyValue && dst.CanSet() && src.Type().AssignableTo(dst.Type()) {
				dst.Set(src)
			}
			break
		}

		if src.Kind() != reflect.Interface {
			if dst.IsNil() || (src.Kind() != reflect.Ptr && overwrite) {
				if dst.CanSet() && (overwrite || isEmptyValue(dst)) {
					dst.Set(src)
				}
			} else if src.Kind() == reflect.Ptr {
				if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
					return
				}
			} else if dst.Elem().Type() == src.Type() {
				if err = deepMerge(dst.Elem(), src, visited, depth+1, config); err != nil {
					return
				}
			} else {
				return ErrDifferentArgumentsTypes
			}
			break
		}

		if dst.IsNil() || overwrite {
			if dst.CanSet() && (overwrite || isEmptyValue(dst)) {
				dst.Set(src)
			}
			break
		}

		if dst.Elem().Kind() == src.Elem().Kind() {
			if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
				return
			}
			break
		}
	default:
		mustSet := (isEmptyValue(dst) || overwrite) && (!isEmptyValue(src) || overwriteWithEmptyValue)
		if mustSet {
			if dst.CanSet() {
				dst.Set(src)
			} else {
				dst = src
			}
		}
	}

	return
}

func merge(dst, src interface{}, opts ...Option) error {
	if dst != nil && reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return ErrNonPointerAgument
	}
	var (
		vDst, vSrc reflect.Value
		err        error
	)

	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	if vDst, vSrc, err = resolveValues(dst, src); err != nil {
		return err
	}
	if vDst.Type() != vSrc.Type() {
		return ErrDifferentArgumentsTypes
	}
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, options)
}

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...Option) error {
	return merge(dst, src, opts...)
}
