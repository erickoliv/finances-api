package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_TableName(t *testing.T) {
	tags := Tag{}
	assert.Equal(t, "public.tags", tags.TableName())
}
