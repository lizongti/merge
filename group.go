package merge

type Group struct {
	sources []any
	opts    []Option
}

func NewGroup(dst any, src any, extra ...any) *Group {
	return &Group{
		sources: append([]any{dst, src}, extra...),
	}
}

func (g *Group) Add(src ...any) *Group {
	g.sources = append(g.sources, src...)
	return g
}

func (g *Group) Opt(opts ...Option) *Group {
	g.opts = append(g.opts, opts...)
	return g
}

func (g *Group) Merge(opts ...Option) (any, error) {
	opts = append(g.opts, opts...)
	ret := g.sources[0]

	for index := 1; index < len(g.sources); index++ {
		var err error

		ret, err = Merge(ret, g.sources[index], opts...)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (g *Group) MustMerge(opts ...Option) any {
	ret, err := g.Merge(opts...)
	if err != nil {
		panic(err)
	}

	return ret
}

func (g *Group) SafeMerge(opts ...Option) (any, error) {
	opts = append(g.opts, opts...)
	ret := g.sources[0]

	for index := 1; index < len(g.sources); index++ {
		var err error

		ret, err = SafeMerge(ret, g.sources[index], opts...)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
