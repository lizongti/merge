package merge

import "fmt"

type SliceStrategy int

const (
	SliceStrategyIgnore SliceStrategy = iota
	SliceStrategyAppend
	SliceStrategyReplaceSlice
	SliceStrategyReplaceElementsDynamic
	SliceStrategyReplaceElementsStatic
	SliceStrategyReplaceDeepDynamic
	SliceStrategyReplaceDeepStatic
)

var sliceStrategyNames = map[SliceStrategy]string{
	SliceStrategyIgnore:                 "Ignore",
	SliceStrategyAppend:                 "Append",
	SliceStrategyReplaceSlice:           "ReplaceSlice",
	SliceStrategyReplaceElementsDynamic: "ReplaceElementsDynamic",
	SliceStrategyReplaceElementsStatic:  "ReplaceElementsStatic",
	SliceStrategyReplaceDeepDynamic:     "ReplaceDeepDynamic",
	SliceStrategyReplaceDeepStatic:      "ReplaceDeepStatic",
}

func (s SliceStrategy) String() string {
	if v, ok := sliceStrategyNames[s]; ok {
		return v
	}
	return fmt.Errorf("%w: %d", ErrInvalidStrategy, s).Error()
}

type StructStrategy int

const (
	StructStrategyIgnore StructStrategy = iota
	StructStrategyReplaceStruct
	StructStrategyReplaceFields
	StructStrategyReplaceDeep
)

var structStrategyNames = map[StructStrategy]string{
	StructStrategyIgnore:        "Ignore",
	StructStrategyReplaceStruct: "ReplaceStruct",
	StructStrategyReplaceFields: "ReplaceFields",
	StructStrategyReplaceDeep:   "ReplaceDeep",
}

func (s StructStrategy) String() string {
	if v, ok := structStrategyNames[s]; ok {
		return v
	}
	return fmt.Errorf("%w: %d", ErrInvalidStrategy, s).Error()
}

type ArrayStrategy int

const (
	ArrayStrategyIgnore ArrayStrategy = iota
	ArrayStrategyAppend
	ArrayStrategyReplaceArray
	ArrayStrategyReplaceElementsDynamic
	ArrayStrategyReplaceElementsStatic
	ArrayStrategyReplaceDeepDynamic
	ArrayStrategyReplaceDeepStatic
)

var arrayStrategyNames = map[ArrayStrategy]string{
	ArrayStrategyIgnore:                 "Ignore",
	ArrayStrategyAppend:                 "Append",
	ArrayStrategyReplaceArray:           "ReplaceArray",
	ArrayStrategyReplaceElementsDynamic: "ReplaceElementsDynamic",
	ArrayStrategyReplaceElementsStatic:  "ReplaceElementsStatic",
	ArrayStrategyReplaceDeepDynamic:     "ReplaceDeepDynamic",
	ArrayStrategyReplaceDeepStatic:      "ReplaceDeepStatic",
}

func (s ArrayStrategy) String() string {
	if v, ok := arrayStrategyNames[s]; ok {
		return v
	}
	return fmt.Errorf("%w: %d", ErrInvalidStrategy, s).Error()
}
