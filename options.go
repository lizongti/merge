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
	arrayResolver   Resolver
	structResolver  Resolver
	sliceResolver   Resolver
	chanResolver    Resolver
	mapResolver     Resolver

	arrayStrategy  ArrayStrategy
	structStrategy StructStrategy
	sliceStrategy  SliceStrategy
	chanStrategy   ChanStrategy
	mapStrategy    MapStrategy
}

var optionsDefault = Options{
	conditions: newConditions(),

	defaultResolver: ResolverNone,
	arrayResolver:   ResolverNone,
	structResolver:  ResolverNone,
	sliceResolver:   ResolverNone,
	chanResolver:    ResolverNone,
	mapResolver:     ResolverNone,

	arrayStrategy:  ArrayStrategyIgnore,
	structStrategy: StructStrategyIgnore,
	sliceStrategy:  SliceStrategyIgnore,
	chanStrategy:   ChanStrategyIgnore,
	mapStrategy:    MapStrategyIgnore,
}

func newOptions(opts []Option) *Options {
	config := optionsDefault
	for _, opt := range opts {
		opt(&config)
	}
	return &config
}

func WithCondition(canCover Condition) Option {
	return func(config *Options) {
		config.conditions = append(config.conditions, canCover)
	}
}

func WithDefaultResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.defaultResolver = resolver
	}
}

func WithArrayResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.arrayResolver = resolver
	}
}

func WithStructResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.structResolver = resolver
	}
}

func WithSliceResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.sliceResolver = resolver
	}
}

func WithChanResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.chanResolver = resolver
	}
}

func WithArrayStrategy(strategy ArrayStrategy) Option {
	return func(config *Options) {
		config.arrayStrategy = strategy
	}
}

func WithStructStrategy(strategy StructStrategy) Option {
	return func(config *Options) {
		config.structStrategy = strategy
	}
}

func WithChanStrategy(strategy ChanStrategy) Option {
	return func(config *Options) {
		config.chanStrategy = strategy
	}
}

func WithSliceStrategy(strategy SliceStrategy) Option {
	return func(config *Options) {
		config.sliceStrategy = strategy
	}
}

func WithMapStrategy(strategy MapStrategy) Option {
	return func(config *Options) {
		config.mapStrategy = strategy
	}
}
