package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"./server"
	"./store"
)

func main() {
	portPtr := flag.String("port", "8080", "Server port. Default: 8080")
	flag.Parse()

	var err error
	logger := log.New(os.Stdout, "faceit-test-commitment", log.LstdFlags|log.Lshortfile)

	logger.Println("preparing SQLite database")
	db, err := store.NewSQLite()
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("SQLite database prepared")

	h := server.NewHandlers(logger, db)
	router := h.SetupRouts()
	srv := server.New(router, fmt.Sprintf(":%s",*portPtr))

	logger.Printf("start to serve localhost at port %s\n", *portPtr)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start %v\n", err)
	}
}
