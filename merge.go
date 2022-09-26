package merge

import "reflect"

type merger struct {
	*Options
	visited map[uintptr]*visit
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

func newMerger(options *Options) *merger {
	return &merger{
		Options: options,
		visited: make(map[uintptr]*visit),
	}
}

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...Option) error {
	return merge(dst, src, opts...)
}

func merge(dst, src interface{}, opts ...Option) error {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)

	if !vDst.CanAddr() {
		return ErrNotAdrressable
	}

	vDst = vDst.Elem()

	return newMerger(options).merge(vDst, vSrc)
}

func (m *merger) merge(dst reflect.Value, src reflect.Value) error {
	var err error

	if dst, err = m.resolve(dst); err != nil {
		return err
	}

	if src, err = m.resolve(src); err != nil {
		return err
	}

	if dst.Kind() != src.Kind() && m.mappingEnabled {
		return m.mapping(dst, src)
	}

	return m.cover(dst, src)
}

func (m *merger) resolve()
