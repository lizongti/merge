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
	resolver      Resolver
	conditions    Conditions
	sliceStrategy SliceStrategy
}

func WithResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.resolver = resolver
	}
}

func WithCondition(canCover Condition) Option {
	return func(config *Options) {
		config.conditions = append(config.conditions, canCover)
	}
}

func WithSliceStrategy(strategy SliceStrategy) Option {
	return func(config *Options) {
		config.sliceStrategy = strategy
	}
}
