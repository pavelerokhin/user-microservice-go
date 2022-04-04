package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"../store"
)

type Handlers struct {
	logger *log.Logger
	db     *store.DB
}

const message = "start here"

func (h *Handlers) AddOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) AllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	id := getIdFromVars(r)

	w.WriteHeader(http.StatusOK)
	fmt.Printf("DELETE user id %s\n", id)
	w.Write([]byte(message + " " + id))
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	id := getIdFromVars(r)

	w.WriteHeader(http.StatusOK)
	fmt.Printf("GET user id %s\n", id)
	w.Write([]byte(message + " " + id))
}

func (h *Handlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	id := getIdFromVars(r)

	w.WriteHeader(http.StatusOK)
	fmt.Printf("UPDATE user id %s\n", id)
	w.Write([]byte(message + " " + id))
}

func NewHandlers(logger *log.Logger, db *store.DB) *Handlers {
	return &Handlers{
		logger: logger,
		db:     db,
	}
}

func getIdFromVars(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}

func (h *Handlers) SetupRouts() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/{filter}", h.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/user/", h.AddOrUpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", h.DeleteUser).Methods(http.MethodDelete)

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
