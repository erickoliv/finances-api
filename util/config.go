package util

import (
	"github.com/ericktm/olivsoft-golang-api/api"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"os"
	"time"
)

// Config is the Application Shared/Singleton props resource
type Config struct {
	ApplicationName string
	DB              *gorm.DB
	Router          *mux.Router
	StartupTime     time.Time
}

// GetConfig is the function designed
// to prepare and return all shared/singleton application props
func GetConfig() Config {
	dbUrl := getEnvConfig("DB_URL")
	_ = getEnvConfig("APP_TOKEN")

	log.Println("url",dbUrl)

	db, err := gorm.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(true)

	// database migrations
	log.Println("start database migrations")
	db.AutoMigrate(&model.Tag{})
	log.Println("stop database migrations")

	r := mux.NewRouter()
	r.Use(authMiddleware)
	r.Use(loggingMiddleware)

	r.Handle("/", api.IndexHandler(db)).Methods("GET")
	r.Handle("/api/tags/{uuid}", api.GetTag(db)).Methods("GET")
	r.Handle("/api/tags/{uuid}", api.UpdateTag(db)).Methods("PUT")
	r.Handle("/api/tags/{uuid}", api.DeleteTag(db)).Methods("DELETE")
	r.Handle("/api/tags", api.GetTags(db)).Methods("GET")
	r.Handle("/api/tags", api.CreateTag(db)).Methods("POST")
	//404 handler
	r.NotFoundHandler = api.PageNotFound()

	cfg := Config{
		"OlivSoft",
		db,
		r,
		time.Now(),
	}

	return cfg
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: cheack token with user token based auth
		if r.Header.Get("Authorization") == os.Getenv("APP_TOKEN") {
			next.ServeHTTP(w, r)
		} else {
			api.UnauthorizedResponse(w)
		}
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("request", r.RemoteAddr, r.Host,r.Method,r.TLS, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func getEnvConfig(s string) string {
	if value, found := os.LookupEnv(s); found {
		return value
	} else {
		log.Fatalf("Environment variable %s not found", s)
	}
	return ""
}