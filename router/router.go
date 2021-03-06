// pkg represents an abstraction on routers of the microservice
// Mux and Chi routers are available

package router

import "net/http"

type Router interface {
	DELETE(uri string, f func(w http.ResponseWriter, r *http.Request))
	GET(uri string, f func(w http.ResponseWriter, r *http.Request))
	POST(uri string, f func(w http.ResponseWriter, r *http.Request))
	SERVE(port string)
}
