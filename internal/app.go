package internal

import (
	"log"
	"syscall"
	"time"

	"github.com/erickoliv/finances-api/account"
	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/index"
	"github.com/erickoliv/finances-api/internal/db"
	"github.com/erickoliv/finances-api/repository/session"
	"github.com/erickoliv/finances-api/repository/sql"
	"github.com/erickoliv/finances-api/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func buildRouter(conn *gorm.DB) *gin.Engine {

	accountRepo := sql.MakeAccounts(conn)
	authenticator := sql.MakeAuthenticator(conn)
	signer := makeJWTSigner()

	r := gin.Default()
	r.GET("/", index.Handler)

	r.Use(db.Middleware(conn))

	security := r.Group("/auth")
	authHandler := auth.NewHTTPHandler(authenticator, signer)
	authHandler.Router(security)

	api := r.Group("/api")
	api.Use(auth.Middleware(signer))

	accounts := account.NewHTTPHandler(accountRepo)
	accounts.Router(api)

	// rest.Routes(api)

	return r
}

// TODO: use a configuration service
func makeJWTSigner() service.Signer {
	appToken, found := syscall.Getenv("APP_TOKEN")
	if !found {
		log.Fatal("env APP_TOKEN not found")
	}

	key := []byte(appToken)
	ttl := time.Hour
	return session.NewJWTSigner(key, ttl)
}

func Run() error {
	conn := db.Prepare()
	defer conn.Close()

	return buildRouter(conn).Run()
}
