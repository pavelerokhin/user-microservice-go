package router

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type muxRouter struct {
	Logger        *log.Logger
	MuxDispatcher *mux.Router
}

func NewMuxRouter(logger *log.Logger) Router {
	return &muxRouter{Logger: logger, MuxDispatcher: mux.NewRouter()}
}

func (mr *muxRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	mr.MuxDispatcher.HandleFunc(uri, f).Methods("DELETE")
}

func (mr *muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	mr.MuxDispatcher.HandleFunc(uri, f).Methods("GET")
}

func (mr *muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	mr.MuxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (mr *muxRouter) SERVE(port string) {
	mr.Logger.Printf("Mux HTTP server running on posrt %v", port)
	mr.Logger.Fatalln(http.ListenAndServe(port, mr.MuxDispatcher))
}
