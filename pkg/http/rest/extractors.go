package rest

import (
	"errors"

	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExtractUser extract and returns a user UUID from gin context
func ExtractUser(c *gin.Context) (uuid.UUID, error) {
	identifier := c.GetString(domain.LoggedUser)
	if len(identifier) == 0 {
		return uuid.Nil, errors.New("user not present in context")
	}

	return uuid.Parse(identifier)
}

// ExtractUUID extract UUID from URL parameters
func ExtractUUID(c *gin.Context) (uuid.UUID, error) {
	pk, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return uuid.Nil, errors.New("uuid parameter is invalid")
	}
	return pk, nil
}
