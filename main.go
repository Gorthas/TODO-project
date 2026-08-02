package main

import (
	"log"
	"os"

	"github.com/Gorthas/TODO-project/pkg/db"
	"github.com/Gorthas/TODO-project/pkg/server"
)

func main() {
	port := "7540"
	dbFile := "scheduler.db"

	manualDBFile := os.Getenv("TODO_DBFILE")
	if len(manualDBFile) > 0 {
		dbFile = manualDBFile
	}

	manualPort := os.Getenv("TODO_PORT")
	if len(manualPort) > 0 {
		port = manualPort
	}
	addr := ":" + port

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
