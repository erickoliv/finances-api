package url

import (
	"github.com/ericktm/olivsoft-golang-api/api"
	"github.com/ericktm/olivsoft-golang-api/auth"
	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PrepareRouter add description
func PrepareRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/", IndexHandler)

	r.Use(database.Middleware(db))

	security := r.Group("/auth")
	auth.Routes(security)

	rest := r.Group("/api")
	rest.Use(auth.Middleware())
	api.TagRoutes(rest)

	return r
}
