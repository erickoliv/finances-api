package main

import (
	"github.com/erickoliv/finances-api/internal"
	"log"
	"os"
)

func main() {

	if err := internal.Run(); err != nil {
		log.Printf("error running application %v", err)
		os.Exit(1)
	}

}
