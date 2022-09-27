package merge

type SliceStrategy int

const (
	SliceStrategyIgnore SliceStrategy = iota
	SliceStrategyAppend
	SliceStrategyReplaceSlice
	SliceStrategyReplaceElement
	SliceStrategyReplaceDeep
)

type StructStrategy int

const (
	StructStrategyIgnore StructStrategy = iota
	StructStrategyReplaceStruct
	StructStrategyReplaceField
	StructStrategyReplaceExportedField
	StructStrategyReplaceDeepField
	StructStrategyReplaceDeepExportedField
)
