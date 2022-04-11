package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

var (
	chiDispatcher = chi.NewRouter()
)

type chiRouter struct{}

func (*chiRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Delete(uri, f)
}

func (*chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, f)
}

func (*chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, f)
}

func (*chiRouter) SERVE(port string) {
	log.Printf("CHI HTTP server running on posrt %v", port)
	log.Fatalln(http.ListenAndServe(port, chiDispatcher))
}

func NewChiRouter() Router {
	return &chiRouter{}
}
