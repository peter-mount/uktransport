package main

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/uktransport/dbrest"
	"log"
)

func main() {
	err := kernel.Launch(&dbrest.DBRest{})
	if err != nil {
		log.Fatal(err)
	}
}
