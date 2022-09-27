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
		merge.WithResolver(merge.ResolverNone)).(*int))

	assert.Equal(t, 1, *merge.MustMerge(&a, &b,
		merge.WithResolver(merge.ResolverBoth)).(*int))

	assert.Equal(t, 1, **merge.MustMerge(&pa, &pb,
		merge.WithResolver(merge.ResolverDeepBoth)).(**int))

	assert.Equal(t, 1, *merge.MustMerge(&a, b,
		merge.WithResolver(merge.ResolverSingle)).(*int))

	assert.Equal(t, 1, *merge.MustMerge(&a, &pb,
		merge.WithResolver(merge.ResolverDeepSingle)).(*int))

	assert.Equal(t, 1, merge.MustMerge(a, &pb,
		merge.WithResolver(merge.ResolverDeepSingle)).(int))

	assert.Equal(t, 10, *merge.MustMerge(&a, b,
		merge.WithResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionTypeCheck)).(*int))
}

func TestFunction(t *testing.T) {
	type f = func() int
	var (
		a = func() int { return 10 }
		b = func() int { return 1 }
	)

	assert.Equal(t, 1, merge.MustMerge(a, b).(f)())

	assert.Equal(t, 1, merge.MustMerge(a, &b,
		merge.WithResolver(merge.ResolverSingle)).(f)())
}

func TestSlice(t *testing.T) {
	type s = []int
	type ss = [][]int
	var (
		a = s{1, 2, 3}
		b = s{4, 5}
		c = ss{{6, 7}, {8, 9}}
		d = ss{{10, 11}}
		e = ss{{12, 13, 14}}
		f = ss{{15}}
	)

	assert.Equal(t, s{1, 2, 3, 4, 5}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyAppend),
	).(s))

	assert.Equal(t, s{4, 5}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceSlice),
	).(s))

	assert.Equal(t, s{4, 5, 3}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElement),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(s))

	assert.Equal(t, s{4, 5, 3}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(s))

	assert.Equal(t, ss{{6, 7}, {8, 9}, {10, 11}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyAppend),
	).(ss))

	assert.Equal(t, ss{{10, 11}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceSlice),
	).(ss))

	assert.Equal(t, ss{{10, 11}, {8, 9}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElement),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(ss))

	assert.Equal(t, ss{{10, 11}, {8, 9}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElement),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(ss))

	assert.Equal(t, ss{{15}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElement),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(ss))

	assert.Equal(t, ss{{15, 7}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeep),
		merge.WithCondition(merge.ConditionSrcIsValid),
	).(ss))
}
