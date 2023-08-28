package merge

import (
	"errors"
	"reflect"
	"unicode/utf8"
)

// Errors reported by merge when it finds invalid arguments.
var (
	ErrNilArguments                = errors.New("src and dst must not be nil")
	ErrDifferentArgumentsTypes     = errors.New("src and dst must be of same type")
	ErrNotSupported                = errors.New("only structs, maps, and slices are supported")
	ErrExpectedMapAsDestination    = errors.New("dst was expected to be a map")
	ErrExpectedStructAsDestination = errors.New("dst was expected to be a struct")
	ErrNonPointerAgument           = errors.New("dst must be a pointer")
)

func hasMergeableFields(dst reflect.Value) (exported bool) {
	for i, n := 0, dst.NumField(); i < n; i++ {
		field := dst.Type().Field(i)
		if field.Anonymous && dst.Field(i).Kind() == reflect.Struct {
			exported = exported || hasMergeableFields(dst.Field(i))
		} else if isExportedComponent(&field) {
			exported = exported || len(field.PkgPath) == 0
		}
	}
	return
}

func isExportedComponent(field *reflect.StructField) bool {
	pkgPath := field.PkgPath
	if len(pkgPath) > 0 {
		return false
	}
	c := field.Name[0]
	if 'a' <= c && c <= 'z' || c == '_' {
		return false
	}
	return true
}

func isExported(field reflect.StructField) bool {
	r, _ := utf8.DecodeRuneInString(field.Name)
	return r >= 'A' && r <= 'Z'
}

// IsReflectNil is the reflect value provided nil
func isReflectNil(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Interface, reflect.Slice, reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr:
		// Both interface and slice are nil if first word is 0.
		// Both are always bigger than a word; assume flagIndir.
		return v.IsNil()
	default:
		return false
	}
}

// During deepMerge, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited are stored in a map indexed by 17 * a1 + a2;
type visit struct {
	typ  reflect.Type
	next *visit
	ptr  uintptr
}

func resolveValues(dst, src interface{}) (vDst, vSrc reflect.Value, err error) {
	if dst == nil || src == nil {
		err = ErrNilArguments
		return
	}
	vDst = reflect.ValueOf(dst).Elem()
	if vDst.Kind() != reflect.Struct && vDst.Kind() != reflect.Map && vDst.Kind() != reflect.Slice {
		err = ErrNotSupported
		return
	}
	vSrc = reflect.ValueOf(src)
	// We check if vSrc is a pointer to dereference it.
	if vSrc.Kind() == reflect.Ptr {
		vSrc = vSrc.Elem()
	}
	return
}

func changeInitialCase(s string, mapper func(rune) rune) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(mapper(r)) + s[n:]
}

func isValidValue(v reflect.Value) bool {
	return v.IsValid()
}

func isZeroValue(v reflect.Value) bool {
	return v.IsZero()
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}

type Range int

const (
	RangeAll Range = iota
	RangeMap
	RangeSlice
	RangeStruct
)

type Level int

const (
	LevelInvalid Level = iota
	LevelZero
	LevelEmpty
	LevelRelevant
)

type Style int

const (
	StyleAll Style = iota
	StyleEach
	StyleRecursive
	StyleAppend
)

type TypeIgoreType int

const (
	TypeIgnoreAll TypeIgoreType = iota
	TypeIgnoreEmpty
	TypeIgnoreZero
	TypeIgnoreNone
)

func ResolveLevel(v reflect.Value) Level {
	switch {
	case !v.IsValid():
		return LevelInvalid
	case v.IsZero():
		return LevelZero
	case (v.Kind() == reflect.Map || v.Kind() == reflect.Slice) && v.Len() == 0:
		return LevelEmpty
	default:
		return LevelRelevant
	}
}
