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
	var a int = 10
	var b int = 1
	ret := merge.MustMerge(&a, &b,
		merge.WithResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionCoverAll))
	assert.Equal(t, 1, *ret.(*int))
}
