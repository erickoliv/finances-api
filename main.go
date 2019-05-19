package main

import (
	"log"

	"github.com/ericktm/olivsoft-golang-api/util"
)

func main() {
	log.Println("Application Start")
	defer log.Println("Application Stop")

	app := util.GetConfig()

	log.Fatal(app.Router.Run())

	defer app.DB.Close()
}
