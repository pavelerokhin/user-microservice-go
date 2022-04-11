package server

import (
	//"crypto/tls"
	"github.com/pavelerokhin/user-microservice-go/router"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server interface {
	SetupRouts(router *router.Router)
}

type serv struct {
	S      *http.Server
	Logger *log.Logger
}

func New(mux *mux.Router, serviceAddr string, logger *log.Logger) *Server {
	//tlsConfig := &tls.Config{
	//	PreferServerCipherSuites: true,
	//	CurvePreferences: []tls.CurveID{
	//		tls.CurveP256,
	//		tls.X25519, //TODO: control this!
	//	},
	//
	//	MinVersion: tls.VersionTLS12,
	//	CipherSuites: []uint16{
	//		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	//		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, // TODO: LOOK SPEC OF RSA AND ECDSA
	//		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	//		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	//		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	//		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	//	},
	//}

	srv := &http.Server{
		Addr:         serviceAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		//TLSConfig:    tlsConfig,
		Handler: mux,
	}

	s := serv{S: srv, Logger: logger}
	return s
}

func (s *serv) SetupRouts(router *router.Router) {
	s.Logger.Println("setting router and handle function")

	router.GET("/users", h.GetUsers)                                  // without pagination
	router.GET("/users/{page-size:[0-9]+}/{page:[0-9]+}", h.GetUsers) // with pagination
	router.POST("/user", h.AddUser)
	router.POST("/user/{id}", h.UpdateUser)
	router.POST("/user/{id}", h.GetUser)
	router.DELETE("/user/{id}", h.DeleteUser)

	router.GET("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

}
