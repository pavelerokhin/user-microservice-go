package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct {
	Logger        *log.Logger
	MuxDispatcher *mux.Router
}

func NewMuxRouter(logger *log.Logger) Router {
	return &muxRouter{Logger: logger, MuxDispatcher: mux.NewRouter()}
}

func (mr *muxRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	mr.MuxDispatcher.HandleFunc(uri, f).Methods(http.MethodDelete)
}

func (mr *muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	mr.MuxDispatcher.HandleFunc(uri, f).Methods(http.MethodGet)
}

func (mr *muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	mr.MuxDispatcher.HandleFunc(uri, f).Methods(http.MethodPost)
}

func (mr *muxRouter) SERVE(port string) {
	mr.Logger.Printf("Mux HTTP server running on port %v", port)
	mr.Logger.Fatalln(http.ListenAndServe(port, mr.MuxDispatcher))
}
