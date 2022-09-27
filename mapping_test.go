package merge_test

import (
	"testing"

	"github.com/cloudlibraries/merge"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert.Equal(t, 1, merge.MustMerge(10, 1,
		merge.WithCondition(merge.ConditionCoverAll)))

	assert.Equal(t, 1, merge.MustMerge(0, 1,
		merge.WithCondition(merge.ConditionCoverZero)))

	assert.Equal(t, 10, merge.MustMerge(10, 1,
		merge.WithCondition(merge.ConditionCoverZero)))
}

func TestPointer(t *testing.T) {
	var (
		a  int  = 10
		b  int  = 1
		pa *int = &a
		pb *int = &b
	)

	assert.Equal(t, &b, merge.MustMerge(&a, &b,
		merge.WithResolver(merge.ResolverNone),
		merge.WithCondition(merge.ConditionCoverAll)).(*int))

	assert.Equal(t, 1, *merge.MustMerge(&a, &b,
		merge.WithResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionCoverAll)).(*int))

	assert.Equal(t, 1, **merge.MustMerge(&pa, &pb,
		merge.WithResolver(merge.ResolverDeepBoth),
		merge.WithCondition(merge.ConditionCoverAll)).(**int))

	assert.Equal(t, 1, *merge.MustMerge(&a, b,
		merge.WithResolver(merge.ResolverSingle),
		merge.WithCondition(merge.ConditionCoverAll)).(*int))

	assert.Equal(t, 1, *merge.MustMerge(&a, &pb,
		merge.WithResolver(merge.ResolverDeepSingle),
		merge.WithCondition(merge.ConditionCoverAll)).(*int))

	assert.Equal(t, 1, merge.MustMerge(a, &pb,
		merge.WithResolver(merge.ResolverDeepSingle),
		merge.WithCondition(merge.ConditionCoverAll)).(int))
}
