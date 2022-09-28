package merge_test

import (
	"testing"

	"github.com/cloudlibraries/merge"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert.Equal(t, 1, merge.MustMerge(10, 1))

	assert.Equal(t, 1, merge.MustMerge(0, 1,
		merge.WithCondition(merge.ConditionDstIsZero)))

	assert.Equal(t, 10, merge.MustMerge(10, 1,
		merge.WithCondition(merge.ConditionDstIsZero)))
}

func TestPointer(t *testing.T) {
	var (
		a  int  = 10
		b  int  = 1
		pa *int = &a
		pb *int = &b
	)

	assert.Equal(t, &b, merge.MustMerge(&a, &b,
		merge.WithDefaultResolver(merge.ResolverNone),
	).(*int))

	assert.Equal(t, 1, *merge.MustMerge(&a, &b,
		merge.WithDefaultResolver(merge.ResolverBoth),
	).(*int))

	assert.Equal(t, 1, **merge.MustMerge(&pa, &pb,
		merge.WithDefaultResolver(merge.ResolverDeepBoth),
	).(**int))

	assert.Equal(t, 1, *merge.MustMerge(&a, b,
		merge.WithDefaultResolver(merge.ResolverSingle),
	).(*int))

	assert.Equal(t, 1, *merge.MustMerge(&a, &pb,
		merge.WithDefaultResolver(merge.ResolverDeepSingle),
	).(*int))

	assert.Equal(t, 1, merge.MustMerge(a, &pb,
		merge.WithDefaultResolver(merge.ResolverDeepSingle),
	).(int))

	assert.Equal(t, 10, *merge.MustMerge(&a, b,
		merge.WithDefaultResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionTypeCheck),
	).(*int))
}

func TestFunction(t *testing.T) {
	type f = func() int
	var (
		a = func() int { return 10 }
		b = func() int { return 1 }
	)

	assert.Equal(t, 1, merge.MustMerge(a, b).(f)())

	assert.Equal(t, 1, merge.MustMerge(a, &b,
		merge.WithDefaultResolver(merge.ResolverSingle)).(f)())
}

func TestSlice(t *testing.T) {
	type s = []int
	type ss = [][]int
	type sps = []*[]int
	var (
		a = s{1, 2, 3}
		b = s{4, 5}
		c = ss{{6, 7}, {8, 9}}
		d = ss{{10, 11}}
		e = ss{{12, 13, 14}}
		f = ss{{15}}
		g = sps{{15}}
	)

	assert.Equal(t, s{1, 2, 3, 4, 5}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyAppend),
	).(s))

	assert.Equal(t, s{4, 5}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceSlice),
	).(s))

	assert.Equal(t, s{4, 5, 3}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElements),
	).(s))

	assert.Equal(t, s{4, 5, 3}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
	).(s))

	assert.Equal(t, ss{{6, 7}, {8, 9}, {10, 11}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyAppend),
	).(ss))

	assert.Equal(t, ss{{10, 11}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceSlice),
	).(ss))

	assert.Equal(t, ss{{10, 11}, {8, 9}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElements),
	).(ss))

	assert.Equal(t, ss{{10, 11}, {8, 9}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElements),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
	).(ss))

	assert.Equal(t, ss{{15}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElements),
	).(ss))

	assert.Equal(t, ss{{15, 7}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
	).(ss))

	assert.Equal(t, ss{{15}, {8, 9}}, merge.MustMerge(c, g,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElements),
		merge.WithDefaultResolver(merge.ResolverSingle),
	).(ss))

	assert.Equal(t, ss{{15, 7}, {8, 9}}, merge.MustMerge(c, g,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
		merge.WithDefaultResolver(merge.ResolverSingle),
	).(ss))
}

func TestStruct(t *testing.T) {
	type a struct {
		A int
	}
	type b struct {
		B int
		a a
	}
	type c struct {
		C int
		a *a
	}
	var (
		// s1 = b{B: 1}
		// s2 = b{a: a{A: 1}}
		s3 = c{C: 1}
		s4 = &c{a: &a{A: 1}}
	)

	// assert.Equal(t, b{B: 1}, merge.MustMerge(s1, s2,
	// 	merge.WithStructStrategy(merge.StructStrategyIgnore),
	// ).(b))

	// assert.Equal(t, b{a: a{A: 1}}, merge.MustMerge(s1, s2,
	// 	merge.WithStructStrategy(merge.StructStrategyReplaceStruct),
	// ).(b))

	// assert.Equal(t, b{B: 0, a: a{A: 1}}, merge.MustMerge(s1, s2,
	// 	merge.WithStructStrategy(merge.StructStrategyReplaceFields),
	// ).(b))

	// assert.Equal(t, b{B: 1, a: a{A: 1}}, merge.MustMerge(s1, s2,
	// 	merge.WithStructStrategy(merge.StructStrategyReplaceFields),
	// 	merge.WithCondition(merge.ConditionSrcIsNotZero),
	// ).(b))

	// assert.Equal(t, b{B: 1, a: a{A: 1}}, merge.MustMerge(s1, s2,
	// 	merge.WithStructStrategy(merge.StructStrategyReplaceDeep),
	// 	merge.WithCondition(merge.ConditionSrcIsNotZero),
	// ).(b))

	// assert.Equal(t, c{a: &a{A: 1}}, merge.MustMerge(s3, s4,
	// 	merge.WithStructStrategy(merge.StructStrategyReplaceStruct),
	// 	merge.WithResolver(merge.ResolverSingle),
	// ).(c))

	// assert.Equal(t, c{C: 1, a: &a{A: 1}}, merge.MustMerge(s3, s4,
	// 	merge.WithStructStrategy(merge.StructStrategyReplaceFields),
	// 	merge.WithResolver(merge.ResolverSingle),
	// 	merge.WithCondition(merge.ConditionSrcIsNotZero),
	// ).(c))

	assert.Equal(t, c{C: 1, a: &a{A: 1}}, merge.MustMerge(s3, s4,
		merge.WithStructStrategy(merge.StructStrategyReplaceDeep),
		merge.WithDefaultResolver(merge.ResolverSingle),
		merge.WithStructResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(c))
}
