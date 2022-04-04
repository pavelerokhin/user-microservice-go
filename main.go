package main

import (
	"log"
	"os"

	"./server"
)


var (
	serviceAddr = ":8080"
)

func main() {
	var err error
	logger := log.New(os.Stdout, "faceit-test-commitment", log.LstdFlags | log.Lshortfile)

	//db, err := sqlx.Open("sqlite", "users.db")
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//
	//err = db.Ping()
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//

	h := server.NewHandlers(logger, nil)
	router := h.SetupRouts()

	srv := server.New(router, serviceAddr)

	logger.Println("server starting")
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start %v", err)
	}
}



