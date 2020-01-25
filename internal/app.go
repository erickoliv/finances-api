package internal

import (
	"github.com/erickoliv/finances-api/internal/db"
	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/index"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func buildRouter(conn *gorm.DB) *gin.Engine {
	repo := sql.MakeAccounts(conn)

	r := gin.Default()
	r.GET("/", index.Handler)

	r.Use(db.Middleware(conn))

	security := r.Group("/auth")
	auth.Routes(security)

	api := r.Group("/api")
	api.Use(auth.Middleware())

	accounts := rest.MakeAccountView(repo)
	accounts.Router(api)

	// rest.Routes(api)

	return r
}

func Run() error {
	conn := db.Prepare()
	defer conn.Close()

	router := buildRouter(conn)

	return router.Run()
}
