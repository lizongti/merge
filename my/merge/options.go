package merge

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrUnknownRange = errors.New("unknown range")
)

type Transformers interface {
	Transformer(reflect.Type) func(dst, src reflect.Value) error
}

type Option func(*Options)

type strategy struct {
	style   Style
	isCover func(dst reflect.Value, src reflect.Value) bool
}

type Options struct {
	transformers Transformers

	overwrite                    bool
	overwriteWithEmptyValue      bool
	overwriteSliceWithEmptyValue bool
	appendSlice                  bool
	typeCheck                    bool
	overwriteRecursively         bool

	Strategies map[Range]strategy
}

// WithTransformers adds transformers to merge, allowing to customize the merging of some types.
func WithTransformers(transformers Transformers) Option {
	return func(config *Options) {
		config.transformers = transformers
	}
}

// WithOverwrite will make merge overwrite non-empty dst attributes with non-empty src attributes values.
func WithOverwrite() Option {
	return func(config *Options) {
		config.overwrite = true
	}
}

// WithOverwriteWithEmptyValue will make merge overwrite non empty dst attributes with empty src attributes values.
func WithOverwriteWithEmptyValue() Option {
	return func(config *Options) {
		config.overwrite = true
		config.overwriteWithEmptyValue = true
	}
}

// WithOverwriteSliceWithEmptyValue will make merge overwrite empty dst slice with empty src slice.
func WithOverwriteSliceWithEmptyValue() Option {
	return func(config *Options) {
		config.overwrite = true
		config.overwriteSliceWithEmptyValue = true
	}
}

// WithAppendSlice will make merge append slices instead of overwriting it.
func WithAppendSlice() Option {
	return func(config *Options) {
		config.appendSlice = true
	}
}

func WithOverwriteRecursively() Option {
	return func(config *Options) {
		config.overwriteRecursively = true
	}
}

func WithStrategy(rng Range, style Style, isCover func(dst reflect.Value, src reflect.Value) bool) Option {
	return func(config *Options) {
		switch rng {
		case RangeAll:
			if _, ok := config.Strategies[RangeMap]; !ok {
				config.Strategies[RangeMap] = strategy{style, isCover}
			}
			if _, ok := config.Strategies[RangeSlice]; !ok {
				config.Strategies[RangeSlice] = strategy{style, isCover}
			}
			if _, ok := config.Strategies[RangeStruct]; !ok {
				config.Strategies[RangeStruct] = strategy{style, isCover}
			}
		case RangeMap:
			config.Strategies[RangeMap] = strategy{style, isCover}
		case RangeSlice:
			config.Strategies[RangeSlice] = strategy{style, isCover}
		case RangeStruct:
			config.Strategies[RangeStruct] = strategy{style, isCover}
		default:
			panic(fmt.Errorf("%w: %v", ErrUnknownRange, rng))
		}
	}
}
