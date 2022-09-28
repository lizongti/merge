package merge

type SliceStrategy int

const (
	SliceStrategyIgnore SliceStrategy = iota
	SliceStrategyAppend
	SliceStrategyReplaceSlice
	SliceStrategyReplaceElements
	SliceStrategyReplaceDeep
)

type StructStrategy int

const (
	StructStrategyIgnore StructStrategy = iota
	StructStrategyReplaceStruct
	StructStrategyReplaceFields
	StructStrategyReplaceDeep
)
