package server

import (
	"../store"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


type Handlers struct {
	logger *log.Logger
	db *store.DB
}

const message = "start here"


func (h *Handlers) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := getIdFromVars(r)

	w.WriteHeader(http.StatusOK)
	fmt.Printf("DELETE user id %s\n", id)
	w.Write([]byte(message + " " + id))
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


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

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := getIdFromVars(r)

	w.WriteHeader(http.StatusOK)
	fmt.Printf("UPDATE user id %s\n", id)
	w.Write([]byte(message + " " + id))
}

func NewHandlers(logger *log.Logger, db *store.DB) *Handlers {
	return &Handlers {
		logger: logger,
		db: db,
	}
}

func getIdFromVars(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}

func (h *Handlers) SetupRouts() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/filter", h.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/user/add", h.AddUser).Methods(http.MethodPost)
	router.HandleFunc("/delete/{id}", h.DeleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/update/{id}", h.UpdateUser).Methods(http.MethodPost)

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
