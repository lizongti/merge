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
	conditions Conditions

	defaultResolver Resolver
	sliceResolver   Resolver
	structResolver  Resolver

	sliceStrategy  SliceStrategy
	structStrategy StructStrategy
}

var optionsDefault = Options{
	defaultResolver: ResolverNone,
	structResolver:  ResolverNone,
	conditions:      newConditions(),
	sliceStrategy:   SliceStrategyIgnore,
	structStrategy:  StructStrategyIgnore,
}

func newOptions(opts []Option) *Options {
	config := optionsDefault
	for _, opt := range opts {
		opt(&config)
	}
	return &config
}

func WithDefaultResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.defaultResolver = resolver
	}
}

func WithSliceResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.sliceResolver = resolver
	}
}

func WithStructResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.structResolver = resolver
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

func WithStructStrategy(strategy StructStrategy) Option {
	return func(config *Options) {
		config.structStrategy = strategy
	}
}
