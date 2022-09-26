package merge_test

import (
	"testing"

	"github.com/cloudlibraries/merge"
	"github.com/stretchr/testify/assert"
)

func TestNum(t *testing.T) {
	assert.Equal(t, 1, merge.MustMerge(10, 1))
}
