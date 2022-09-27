package merge

type SliceStrategy int

const (
	SliceStrategyNone SliceStrategy = iota
	SliceStrategyAppend
	SliceStrategyReplaceSlice
	SliceStrategyReplaceElem
	SliceStrategyReplaceDeep
)
