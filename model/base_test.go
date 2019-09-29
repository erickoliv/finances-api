package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseModel_IsNew(t *testing.T) {
	new := BaseModel{}

	assert.True(t, new.IsNew())
}
