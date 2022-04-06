package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"

	"../store"
)

type Handlers struct {
	logger *log.Logger
	db     *store.DB
}


func (h *Handlers) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	user, err := unmarshalUserFromRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.db.DB.Save(user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User has been successfully added"))
}

func (h *Handlers) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users []store.User
	h.db.DB.Find(&users)

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	id := getIdFromVars(r)

	var user store.User
	h.db.DB.Where("id = ?", id).Find(&user)
	h.db.DB.Delete(&user)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("The user with id %s has been deleted successfully", id)))
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	var  err error
	var users []store.User

	var filters *store.User

	filters, err = unmarshalUserFromRequestBody(r)
	emptyBodyErrorImpl := &EmptyBody{}
	if err != nil && !errors.As(err, &emptyBodyErrorImpl) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	vars := mux.Vars(r)
	if vars["pagination-size"] != "" {
		pageSize, err := strconv.Atoi(vars["pagination-size"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Cannot get pagination limit"))
			return
		}

		page, err := strconv.Atoi(vars["page"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Cannot get page number for pagination"))
			return
		}
		if pageSize <= 0 {
			pageSize = 1
		}

		if filters != nil {
			h.db.DB.Scopes(Paginate(page, pageSize)).Where(&filters).Find(&users)
		} else {
			h.db.DB.Scopes(Paginate(page, pageSize)).Find(&users)
		}
	} else {
		// show all users
		if filters != nil {
			h.db.DB.Where(&filters).Find(&users)
		} else {
			h.db.DB.Find(&users)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&users)
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	id := getIdFromVars(r)
	var user store.User
	h.db.DB.Where("id = ?", id).Find(&user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}

func (h *Handlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	id := getIdFromVars(r)
	var user store.User
	h.db.DB.Where("id = ?", id).Find(&user)

	if &user == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Cannot find user with id %s", id)))
		return
	}

	newUser, err := unmarshalUserFromRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mergeUserObjects(&user, newUser)

	h.db.DB.Save(user)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User with id %s has been updated successfully", id)))
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

type EmptyBody struct {
	message string
	Err error
}

func (e *EmptyBody) Error() string {
	return e.Err.Error()
}

func unmarshalUserFromRequestBody(r *http.Request) (*store.User, error) {
	//// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	//// response body. A request body larger than that will now result in
	//// Decode() returning a "http: request body too large" error.
	//r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var user store.User
	err := dec.Decode(&user)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			return nil, fmt.Errorf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, fmt.Errorf("Request body contains badly-formed JSON")

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			return nil, fmt.Errorf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return nil, fmt.Errorf("Request body contains unknown field %s", fieldName)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			return nil, &EmptyBody{"Request body must not be empty", err} //, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return nil, fmt.Errorf("Request body must not be larger than 1MB") //, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return nil, fmt.Errorf("Internal server error") //, http.StatusInternalServerError)
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return nil, fmt.Errorf("Request body must only contain a single JSON object")
	}

	return &user, nil
}

func mergeUserObjects(userSource, userDest *store.User) {
	if userDest.FirstName != "" {
		userSource.FirstName = userDest.FirstName
	}

	if userDest.LastName != "" {
		userSource.LastName = userDest.LastName
	}

	if userDest.Nickname != "" {
		userSource.Nickname = userDest.Nickname
	}

	if userDest.Country != "" {
		userSource.Country = userDest.Country
	}

	if userDest.Email != "" {
		userSource.Email = userDest.Email
	}

	if userDest.Password != "" {
		userSource.Password = userDest.Password
	}
}

func (h *Handlers) SetupRouts() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", h.GetUsers).Methods(http.MethodGet, http.MethodPost) // without pagination
	router.HandleFunc("/users/{pagination-size:[0-9]+}/{page:[0-9]+}", h.GetUsers).Methods(http.MethodGet, http.MethodPost) // with pagination
	router.HandleFunc("/user", h.AddUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.UpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", h.DeleteUser).Methods(http.MethodDelete)

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
