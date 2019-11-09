package main

import (
	"github.com/ericktm/olivsoft-golang-api/internal"
	"log"
	"os"
)

func main() {

	if err := internal.Run(); err != nil {
		log.Printf("error running application %v", err)
		os.Exit(1)
	}

}
