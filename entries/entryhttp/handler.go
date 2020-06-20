package entryhttp

import (
	"net/http"

	"github.com/erickoliv/finances-api/entries"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo entries.Repository
}

// NewHandler receives a Entry sql repository, returning a HTTP entry handler
func NewHandler(repo entries.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (view *Handler) Router(group *gin.RouterGroup) {
	group.GET("/entries", view.GetEntries)
	group.GET("/entries/:uuid", view.GetEntry)
	group.POST("/entries", view.CreateEntry)
	group.PUT("/entries/:uuid", view.UpdateEntry)
	group.DELETE("/entries/:uuid", view.DeleteEntry)
}

// GetEntries return all entries
func (view *Handler) GetEntries(c *gin.Context) {
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

// CreateEntry can be called to create a new instance of Entry on database, using props provided via http request
func (view Handler) CreateEntry(c *gin.Context) {
	row := &entries.Entry{}
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

// GetEntry can be called to get a specific entry from the database. The uuid used to query is part of the url
func (view Handler) GetEntry(c *gin.Context) {
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
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	c.JSON(http.StatusOK, row)
}

// UpdateEntry can be called to update a specific entry. The uuid used to query is part of the url
func (view Handler) UpdateEntry(c *gin.Context) {
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

	newEntry := &entries.Entry{
		UUID:  pk,
		Owner: user,
	}
	if err := c.Bind(newEntry); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	current, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if current == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	if err := view.repo.Save(c, newEntry); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newEntry)
}

// DeleteEntry can be used to logical delete a row entry from the database.
func (view Handler) DeleteEntry(c *gin.Context) {
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
