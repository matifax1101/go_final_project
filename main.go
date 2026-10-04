package main

import (
	"log"
	"net/http"
	"os"

	"github.com/matifax1101/go_final_project/pkg/api"
	"github.com/matifax1101/go_final_project/pkg/db"
)

const defaultPort = "7540"
const defaultDBFile = "scheduler.db"

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Server started at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
