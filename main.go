package main

import (
	"log"
	"net/http"
	"os"

	"./server"
	"./store"
)


var (
	serviceAddr = ":8080"
)

func main() {
	logger := log.New(os.Stdout, "faceit-test-commitment", log.LstdFlags | log.Lshortfile)

	db, err := store.NewSQLite()
	if err != nil {
		logger.Fatalln(err)
	}

	h := server.NewHandlers(logger, db)

	mux := http.NewServeMux()
	h.SetupRouts(mux)

	srv := server.New(mux, serviceAddr)

	logger.Println("server starting")
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start %v", err)
	}
}



