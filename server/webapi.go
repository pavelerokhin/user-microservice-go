package server

import (
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

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
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


func (h *Handlers) SetupRouts(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Logger(h.Home))

}
