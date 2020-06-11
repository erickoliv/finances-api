package categoryhttp

import (
	"net/http"

	"github.com/erickoliv/finances-api/categories"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo categories.Repository
}

// NewHandler receives a Category sql repository, returning a HTTP category handler
func NewHandler(repo categories.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (view *Handler) Router(group *gin.RouterGroup) {
	group.GET("/categories", view.GetCategories)
	group.GET("/categories/:uuid", view.GetCategory)
	group.POST("/categories", view.CreateCategory)
	group.PUT("/categories/:uuid", view.UpdateCategory)
	group.DELETE("/categories/:uuid", view.DeleteCategory)
}

// GetCategories return all categories
func (view *Handler) GetCategories(c *gin.Context) {
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

// CreateCategory can be called to create a new instance of Category on database, using props provided via http request
func (view Handler) CreateCategory(c *gin.Context) {
	row := &categories.Category{}
	if err := c.ShouldBind(row); err != nil {
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

	row.Owner = user

	if err := view.repo.Save(c, row); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, row)
}

// GetCategory can be called to get a specific category from the database. The uuid used to query is part of the url
func (view Handler) GetCategory(c *gin.Context) {
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
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}

	c.JSON(http.StatusOK, row)
}

// UpdateCategory can be called to update a specific category. The uuid used to query is part of the url
func (view Handler) UpdateCategory(c *gin.Context) {
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

	newCategory := &categories.Category{
		UUID:  pk,
		Owner: user,
	}
	if err := c.Bind(newCategory); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	current, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if current == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}

	if err := view.repo.Save(c, newCategory); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newCategory)
}

// DeleteCategory can be used to logical delete a row category from the database.
func (view Handler) DeleteCategory(c *gin.Context) {
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
