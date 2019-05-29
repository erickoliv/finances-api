package main

import (
	"log"

	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/ericktm/olivsoft-golang-api/url"
)

func main() {
	db := database.PrepareDatabase()
	defer db.Close()

	app := url.PrepareRouter(db)
	log.Fatal(app.Run())
}
