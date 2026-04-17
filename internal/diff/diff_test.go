package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare_Added(t *testing.T) {
	prev := map[string]string{}
	curr := map[string]string{"FOO": "bar"}
	changes := Compare(prev, curr)
	assert.Len(t, changes, 1)
	assert.Equal(t, Added, changes[0].Type)
	assert.Equal(t, "FOO", changes[0].Key)
	assert.Equal(t, "bar", changes[0].NewVal)
}

func TestCompare_Removed(t *testing.T) {
	prev := map[string]string{"FOO": "bar"}
	curr := map[string]string{}
	changes := Compare(prev, curr)
	assert.Len(t, changes, 1)
	assert.Equal(t, Removed, changes[0].Type)
	assert.Equal(t, "FOO", changes[0].Key)
	assert.Equal(t, "bar", changes[0].OldVal)
}

func TestCompare_Updated(t *testing.T) {
	prev := map[string]string{"FOO": "old"}
	curr := map[string]string{"FOO": "new"}
	changes := Compare(prev, curr)
	assert.Len(t, changes, 1)
	assert.Equal(t, Updated, changes[0].Type)
	assert.Equal(t, "old", changes[0].OldVal)
	assert.Equal(t, "new", changes[0].NewVal)
}

func TestCompare_NoChanges(t *testing.T) {
	prev := map[string]string{"FOO": "bar"}
	curr := map[string]string{"FOO": "bar"}
	changes := Compare(prev, curr)
	assert.Empty(t, changes)
}

func TestSummary_NoChanges(t *testing.T) {
	assert.Equal(t, "no changes detected", Summary(nil))
}

func TestSummary_WithChanges(t *testing.T) {
	changes := []Change{
		{Type: Added},
		{Type: Added},
		{Type: Updated},
		{Type: Removed},
	}
	result := Summary(changes)
	assert.Equal(t, "2 added, 1 updated, 1 removed", result)
}
