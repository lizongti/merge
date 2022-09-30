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
	type (
		s   = []int
		ss  = [][]int
		sps = []*[]int
	)
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
		merge.WithSliceStrategy(merge.SliceStrategyReplace),
	).(s))

	assert.Equal(t, s{4, 5}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyRefer),
	).(s))

	assert.Equal(t, s{4, 5, 3}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElementsDynamic),
	).(s))

	assert.Equal(t, s{4, 5, 3}, merge.MustMerge(a, b,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeepDynamic),
	).(s))

	assert.Equal(t, ss{{6, 7}, {8, 9}, {10, 11}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyAppend),
	).(ss))

	assert.Equal(t, ss{{10, 11}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplace),
	).(ss))

	assert.Equal(t, ss{{10, 11}, {8, 9}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElementsDynamic),
	).(ss))

	assert.Equal(t, ss{{10, 11}, {8, 9}}, merge.MustMerge(c, d,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeepDynamic),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElementsDynamic),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeepDynamic),
	).(ss))

	assert.Equal(t, ss{{12, 13, 14}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElementsStatic),
	).(ss))

	assert.Equal(t, ss{{12, 13}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeepStatic),
	).(ss))

	assert.Equal(t, ss{{15}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElementsDynamic),
	).(ss))

	assert.Equal(t, ss{{15, 7}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeepDynamic),
	).(ss))

	assert.Equal(t, ss{{15}, {8, 9}}, merge.MustMerge(c, g,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceElementsDynamic),
		merge.WithSliceResolver(merge.ResolverSingle),
	).(ss))

	assert.Equal(t, ss{{15, 7}, {8, 9}}, merge.MustMerge(c, g,
		merge.WithSliceStrategy(merge.SliceStrategyReplaceDeepDynamic),
		merge.WithSliceResolver(merge.ResolverSingle),
	).(ss))
}

func TestStruct(t *testing.T) {
	type (
		a struct {
			A int
		}
		b struct {
			B int
			a a
		}
		c struct {
			C int
			a *a
		}
	)
	var (
		s1 = b{B: 1}
		s2 = b{a: a{A: 1}}
		s3 = c{C: 1}
		s4 = &c{a: &a{A: 1}}
	)

	assert.Equal(t, b{B: 1}, merge.MustMerge(s1, s2,
		merge.WithStructStrategy(merge.StructStrategyIgnore),
	).(b))

	assert.Equal(t, b{a: a{A: 1}}, merge.MustMerge(s1, s2,
		merge.WithStructStrategy(merge.StructStrategyReplace),
	).(b))

	assert.Equal(t, b{B: 0, a: a{A: 1}}, merge.MustMerge(s1, s2,
		merge.WithStructStrategy(merge.StructStrategyReplaceFields),
	).(b))

	assert.Equal(t, b{B: 1, a: a{A: 1}}, merge.MustMerge(s1, s2,
		merge.WithStructStrategy(merge.StructStrategyReplaceFields),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(b))

	assert.Equal(t, b{B: 1, a: a{A: 1}}, merge.MustMerge(s1, s2,
		merge.WithStructStrategy(merge.StructStrategyReplaceDeep),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(b))

	assert.Equal(t, c{a: &a{A: 1}}, merge.MustMerge(s3, s4,
		merge.WithStructStrategy(merge.StructStrategyReplace),
		merge.WithStructResolver(merge.ResolverSingle),
	).(c))

	assert.Equal(t, c{C: 1, a: &a{A: 1}}, merge.MustMerge(s3, s4,
		merge.WithStructStrategy(merge.StructStrategyReplaceFields),
		merge.WithStructResolver(merge.ResolverSingle),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(c))

	assert.Equal(t, c{C: 1, a: &a{A: 1}}, merge.MustMerge(s3, s4,
		merge.WithStructStrategy(merge.StructStrategyReplaceDeep),
		merge.WithDefaultResolver(merge.ResolverSingle),
		merge.WithStructResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(c))
}

func TestArray(t *testing.T) {
	type (
		a2    = [2]int
		a2a2  = [2][2]int
		a2pa2 = [2]*[2]int
	)
	var (
		a = a2{1, 2}
		b = a2{3, 4}
		c = a2a2{{6, 7}, {8, 9}}
		d = a2a2{{10, 11}}
		e = a2a2{{12, 13}}
		f = a2pa2{{15, 16}}
	)

	assert.Equal(t, a2{1, 2}, merge.MustMerge(a, b,
		merge.WithArrayStrategy(merge.ArrayStrategyIgnore),
	).(a2))

	assert.Equal(t, a2{3, 4}, merge.MustMerge(a, b,
		merge.WithArrayStrategy(merge.ArrayStrategyReplace),
	).(a2))

	assert.Equal(t, a2a2{{10, 11}}, merge.MustMerge(c, d,
		merge.WithArrayStrategy(merge.ArrayStrategyReplace),
	).(a2a2))

	assert.Equal(t, a2a2{{12, 13}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithArrayStrategy(merge.ArrayStrategyReplaceElements),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(a2a2))

	assert.Equal(t, a2a2{{12, 13}, {8, 9}}, merge.MustMerge(c, e,
		merge.WithArrayStrategy(merge.ArrayStrategyReplaceDeep),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(a2a2))

	assert.Equal(t, a2a2{{15, 16}, {8, 9}}, merge.MustMerge(c, f,
		merge.WithArrayStrategy(merge.ArrayStrategyReplaceDeep),
		merge.WithArrayResolver(merge.ResolverSingle),
		merge.WithCondition(merge.ConditionSrcIsNotZero),
	).(a2a2))
}

func TestChan(t *testing.T) {
	type (
		ci = chan int
	)
	var (
		s2c = func(s []int) ci {
			c := make(ci, len(s))
			for i := 0; i < len(s); i++ {
				c <- s[i]
			}
			return c
		}
		c2s = func(c ci) []int {
			s := make([]int, 0, cap(c))
			for i, n := 0, len(c); i < n; i++ {
				s = append(s, <-c)
			}
			for i := 0; i < len(s); i++ {
				c <- s[i]
			}
			return s
		}

		s1 = []int{1, 2}
		s2 = []int{3, 4, 5}
		s3 = []int{1, 2, 3, 4, 5}
		c1 = s2c(s1)
		c2 = s2c(s2)
		c3 = s2c(s3)
	)

	assert.Equal(t, c2s(c1), c2s(merge.MustMerge(c1, c2,
		merge.WithChanStrategy(merge.ChanStrategyIgnore),
	).(ci)))

	assert.Equal(t, c2s(c2), c2s(merge.MustMerge(c1, c2,
		merge.WithChanStrategy(merge.ChanStrategyRefer),
	).(ci)))

	assert.Equal(t, c2s(c2), c2s(merge.MustMerge(c1, c2,
		merge.WithChanStrategy(merge.ChanStrategyReplace),
	).(ci)))

	assert.Equal(t, []int{3, 4, 5, 4, 5}, c2s(merge.MustMerge(c3, c2,
		merge.WithChanStrategy(merge.ChanStrategyReplaceElements),
	).(ci)))

	assert.Equal(t, []int{3, 4, 5, 4, 5}, c2s(merge.MustMerge(c3, c2,
		merge.WithChanStrategy(merge.ChanStrategyReplaceDeep),
	).(ci)))
}
