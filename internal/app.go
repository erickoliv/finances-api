package internal

import (
	"github.com/ericktm/olivsoft-golang-api/internal/db"
	"github.com/ericktm/olivsoft-golang-api/pkg/http/auth"
	"github.com/ericktm/olivsoft-golang-api/pkg/http/index"
	"github.com/ericktm/olivsoft-golang-api/pkg/http/rest"
	"github.com/ericktm/olivsoft-golang-api/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func buildRouter(conn *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/", index.Handler)

	r.Use(db.Middleware(conn))

	security := r.Group("/auth")
	auth.Routes(security)

	api := r.Group("/api")
	api.Use(auth.Middleware())

	repo := sql.MakeAccounts(conn)
	accounts := rest.MakeAccountView(repo)
	api.POST("/accounts", accounts.CreateAccount)
	api.GET("/accounts/:uuid", accounts.GetAccount)
	api.PUT("/accounts/:uuid", accounts.UpdateAccount)
	api.DELETE("/accounts/:uuid", accounts.DeleteAccount)
	api.GET("/accounts", accounts.GetAccounts)

	// rest.Routes(api)

	return r
}

func Run() error {
	conn := db.Prepare()
	defer conn.Close()

	router := buildRouter(conn)

	return router.Run()
}
