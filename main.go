package main

import (
	"github.com/ericktm/olivsoft-golang-api/util"
	"log"
	"net/http"
)

func main() {
	log.Println("Application Start")
	defer log.Println("Application Stop")

	app := util.GetConfig()

	log.Fatal(http.ListenAndServe(":8080", app.Router))

	defer app.DB.Close()
}
