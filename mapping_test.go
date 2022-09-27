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
