package main

import (
	"log"

	"github.com/SealTV/cloudygo/kvstore"
)

func main() {
	if err := kvstore.InitializeTransactionLog(); err != nil {
		log.Fatal(err)
	}

	server := kvstore.NewServer()
	log.Fatal(server.Run(":8080"))
}
