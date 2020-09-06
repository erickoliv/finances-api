package internal

import (
	"os"

	"github.com/erickoliv/finances-api/accounts/accounthttp"
	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/auth/authhttp"
	"github.com/erickoliv/finances-api/categories/categoryhttp"
	"github.com/erickoliv/finances-api/entries/entryhttp"
	"github.com/erickoliv/finances-api/index"
	"github.com/erickoliv/finances-api/internal/cfg"
	"github.com/erickoliv/finances-api/internal/database"
	"github.com/erickoliv/finances-api/repository/session"
	"github.com/erickoliv/finances-api/tags/taghttp"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type (
	App struct {
		signer   auth.SessionSigner
		sqlStore database.SQLStore
		routes   []Router
	}
	Router interface {
		Router(*gin.RouterGroup)
	}
)

func Run() error {
	config, err := cfg.Load(os.Getenv)
	if err != nil {
		return errors.WithStack(err)
	}
	// add config struct validation here

	db, err := database.Connect(&config.DB)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := database.Migrate(db); err != nil {
		return errors.WithStack(err)
	}

	sqlStore := database.BuildSQLStore(db)

	app := App{
		signer:   makeJWTSigner(&config.JWT),
		sqlStore: sqlStore,
		routes:   buildAPIRoutes(sqlStore),
	}

	router := buildRouter(app)

	return router.Run()
}

func buildAPIRoutes(store database.SQLStore) []Router {
	return []Router{
		accounthttp.NewHTTPHandler(store.Account),
		taghttp.NewHTTPHandler(store.Tag),
		entryhttp.NewHandler(store.Entry),
		categoryhttp.NewHandler(store.Category),
	}
}

func buildRouter(app App) *gin.Engine {
	r := gin.Default()
	r.GET("/", index.Handler)

	security := r.Group("/auth")

	authHandler := authhttp.NewHTTPHandler(app.sqlStore.Auth, app.signer)
	authHandler.Router(security)

	api := r.Group("/api")
	api.Use(authhttp.Middleware(app.signer))

	for _, router := range app.routes {
		router.Router(api)
	}

	return r
}

// TODO: use a configuration service
func makeJWTSigner(config *cfg.Auth) auth.SessionSigner {
	key := []byte(config.Token)

	return session.NewJWTSigner(key, config.TTL)
}
