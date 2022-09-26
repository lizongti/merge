package merge

import (
	"reflect"
)

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
	StyleDeep
	StyleAppend
)

type TypeIgoreType int

const (
	TypeIgnoreAll TypeIgoreType = iota
	TypeIgnoreEmpty
	TypeIgnoreZero
	TypeIgnoreNone
)

type Resolver int

const (
	ResolverNone Resolver = iota
	ResolverNormal
	ResolverDeep
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

type Option func(*Options)

type strategy struct {
	style   Style
	isCover func(dst reflect.Value, src reflect.Value) bool
}

type Options struct {
	mappingEnabled bool
	resolver       Resolver
	style          Style
	// Strategies     map[Range]strategy
}

func WithResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.resolver = resolver
	}
}

func WithStyle(style Style) Option {
	return func(config *Options) {

	}
}

// func WithStrategy(rng Range, style Style, isCover func(dst reflect.Value, src reflect.Value) bool) Option {
// 	return func(config *Options) {
// 		switch rng {
// 		case RangeAll:
// 			if _, ok := config.Strategies[RangeMap]; !ok {
// 				config.Strategies[RangeMap] = strategy{style, isCover}
// 			}
// 			if _, ok := config.Strategies[RangeSlice]; !ok {
// 				config.Strategies[RangeSlice] = strategy{style, isCover}
// 			}
// 			if _, ok := config.Strategies[RangeStruct]; !ok {
// 				config.Strategies[RangeStruct] = strategy{style, isCover}
// 			}
// 		case RangeMap:
// 			config.Strategies[RangeMap] = strategy{style, isCover}
// 		case RangeSlice:
// 			config.Strategies[RangeSlice] = strategy{style, isCover}
// 		case RangeStruct:
// 			config.Strategies[RangeStruct] = strategy{style, isCover}
// 		default:
// 			panic(fmt.Errorf("%w: %v", ErrUnknownRange, rng))
// 		}
// 	}
// }
