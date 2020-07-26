package taghttp

import (
	"net/http"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/tags"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// HTTPHandler is a interface to expose tag manipulation using HTTP
type HTTPHandler interface {
	Router(*gin.RouterGroup)
}

type Handler struct {
	repo tags.Repository
}

// NewHTTPHandler receives a Tag Service, returning a internal a HTTP tag handler
func NewHTTPHandler(repo tags.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (view Handler) Router(group *gin.RouterGroup) {
	group.GET("/tags", view.GetTags)
	group.GET("/tags/:uuid", view.GetTag)
	group.POST("/tags", view.CreateTag)
	group.PUT("/tags/:uuid", view.UpdateTag)
	group.DELETE("/tags/:uuid", view.DeleteTag)
}

// GetTags return all tags
func (view Handler) GetTags(c *gin.Context) {
	query, err := rest.ExtractFilters(c, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := view.repo.Query(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := rest.PaginatedMessage{
		Page:  query.Page,
		Data:  &result,
		Count: len(result),
	}

	c.JSON(http.StatusOK, &response)
}

// CreateTag can be called to create a new instance of Tag on database, using props provided via http request
func (view Handler) CreateTag(c *gin.Context) {
	tag := &tags.Tag{}
	if err := c.ShouldBind(tag); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	tag.Owner = user

	if err := view.repo.Save(c, tag); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTag can be called to get a specific tag from the database. The uuid used to query is part of the url
func (view Handler) GetTag(c *gin.Context) {
	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	pk, err := rest.ExtractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	row, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if row.UUID == uuid.Nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "tag not found"})
		return
	}

	c.JSON(http.StatusOK, row)
}

// UpdateTag can be called to update a specific tag. The uuid used to query is part of the url
func (view Handler) UpdateTag(c *gin.Context) {
	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	pk, err := rest.ExtractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newTag := &tags.Tag{
		UUID:  pk,
		Owner: user,
	}
	if err := c.Bind(newTag); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	current, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if current == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "tag not found"})
		return
	}

	if err := view.repo.Save(c, newTag); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTag)
}

// DeleteTag can be used to logical delete a row tag from the database.
func (view Handler) DeleteTag(c *gin.Context) {
	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	pk, err := rest.ExtractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := view.repo.Delete(c, pk, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
