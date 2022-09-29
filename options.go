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
	arrayResolver   Resolver
	chanResolver    Resolver

	sliceStrategy  SliceStrategy
	structStrategy StructStrategy
	arrayStrategy  ArrayStrategy
	chanStrategy   ChanStrategy
}

var optionsDefault = Options{
	conditions: newConditions(),

	defaultResolver: ResolverNone,
	structResolver:  ResolverNone,
	arrayResolver:   ResolverNone,
	chanResolver:    ResolverNone,

	sliceStrategy:  SliceStrategyIgnore,
	structStrategy: StructStrategyIgnore,
	arrayStrategy:  ArrayStrategyIgnore,
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

func WithArrayResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.arrayResolver = resolver
	}
}

func WithChanResolver(resolver Resolver) Option {
	return func(config *Options) {
		config.chanResolver = resolver
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

func WithArrayStrategy(strategy ArrayStrategy) Option {
	return func(config *Options) {
		config.arrayStrategy = strategy
	}
}

func WithChanStrategy(strategy ChanStrategy) Option {
	return func(config *Options) {
		config.chanStrategy = strategy
	}
}
