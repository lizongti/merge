package merge

import "reflect"

func merge(dst, src interface{}, opts ...Option) error {
	var err error

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)

	if vDst, err = resolveAddressableContainer(vDst); err != nil {
		return err
	}

	if vSrc, err = resolveContainer(vSrc); err != nil {
		return err
	}

	return newMerger(options).merge(vDst, vSrc)
}

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...Option) error {
	return merge(dst, src, opts...)
}
