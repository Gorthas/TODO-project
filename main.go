package main

import (
	"log"
	"os"

	"github.com/Gorthas/TODO-project/pkg/server"
)

func main() {
	port := "7540"
	manualPort := os.Getenv("TODO_PORT")
	if len(manualPort) > 0 {
		port = manualPort
	}
	addr := ":" + port

	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
