package merge

type Range int

const (
	RangeAll Range = iota
	RangeMap
	RangeSlice
	RangeStruct
	RangeDefault
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

type Option func(*Options)

type Options struct {
	mappingEnabled bool
	resolver       Resolver

	conditions Conditions
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

func WithCondition(canCover Condition) Option {
	return func(config *Options) {
		config.conditions = append(config.conditions, canCover)
	}
}
