package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)


type Handlers struct {
	logger *log.Logger
	db *sqlx.DB
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
	w.Write([]byte(message + " " + id))
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Println("received request")
		startTime := time.Now()
		defer h.logger.Printf("received processed in %s\n", time.Now().Sub(startTime))
		next(w,r)
	}
}

func NewHandlers(logger *log.Logger, db *sqlx.DB) *Handlers {
	return &Handlers {
		logger: logger,
		db: db,
	}
}

func getIdFromVars(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}

func (h *Handlers) SetupRouts(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Logger(h.GetAllUsers))
	mux.HandleFunc("/get/{id}", h.Logger(h.GetUser))
	mux.HandleFunc("/filter", h.Logger(h.GetAllUsers))
	mux.HandleFunc("/add", h.Logger(h.AddUser))
	mux.HandleFunc("/delete/{id}", h.Logger(h.DeleteUser))
	mux.HandleFunc("/update/{id}", h.Logger(h.UpdateUser))
	mux.HandleFunc("/heartbeat", h.Logger(h.Heartbeat))
}
