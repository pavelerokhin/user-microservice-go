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
	"../utils"

)

type Handlers struct {
	logger *log.Logger
	db     *store.DB
}


func (h *Handlers) AddUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("request add a new user")
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			h.logger.Println(fmt.Sprintf("error: %s", msg))
			return
		}
	}

	user, rerr, statusCode := unmarshalUserFromRequestBody(r)
	if rerr != nil {
		w.WriteHeader(statusCode)
		msg := fmt.Sprintf("error adding user: %s", rerr.Error())
		w.Write([]byte(msg))
		h.logger.Println(msg)
		return
	}

	emptyFields := utils.CheckEmptyFields(user)
	if len(emptyFields) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Request has empty fields: %s", emptyFields)
		h.logger.Printf(msg)
		_, err := w.Write([]byte(msg))

		if err != nil {
			h.logger.Printf("error: User has not been added, and there's a problem returning the response: %v", err)
		}
		return
	}

	h.db.DB.Save(user)
	_, err := w.Write([]byte("User has been successfully added"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("error: User has been successfully added, but there's a problem returning the response: %v", err)

		return
	}

	h.logger.Println("User has been successfully added")
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("request delete user")

	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	id := getIdFromVars(r)

	var user store.User
	h.db.DB.Where("id = ?", id).Find(&user)
	h.db.DB.Delete(&user)

	msg := fmt.Sprintf("User with id %s has been deleted successfully", id)
	_, err := w.Write([]byte(msg))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("error: User has been successfully deleted, but there's a problem returning the response: %v", err)

		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Println(msg)
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("request list users: ")

	var  err error
	var users []store.User

	var filters *store.User

	filters, rerr, statusCode := unmarshalUserFromRequestBody(r)
	emptyBodyErrorImpl := &EmptyBody{}
	if rerr != nil && !errors.As(rerr, &emptyBodyErrorImpl) {
		w.WriteHeader(statusCode)
		h.logger.Println(fmt.Sprintf("error listing users: %s", rerr))
		return
	}

	vars := mux.Vars(r)
	if vars["pagination-size"] != "" {
		h.logger.Println("* with pagination")

		pageSize, err := strconv.Atoi(vars["pagination-size"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Cannot get pagination limit: %v", err)
			h.logger.Println(msg)

			_, err = w.Write([]byte(msg))
			if err != nil {
				h.logger.Println(fmt.Sprintf("error returning the reponse: %v", err))
			}
			return
		}

		page, err := strconv.Atoi(vars["page"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Cannot get page number for pagination: %v", err)
			h.logger.Println(msg)

			_, err = w.Write([]byte(msg))
			if err != nil {
				h.logger.Println(fmt.Sprintf("error returning the reponse: %v", err))
			}
			return
		}
		if pageSize <= 0 {
			pageSize = 1
		}

		if filters != nil {
			h.logger.Println("* with filtering")
			h.db.DB.Scopes(Paginate(page, pageSize)).Where(&filters).Find(&users)
		} else {
			h.logger.Println("* without filtering")
			h.db.DB.Scopes(Paginate(page, pageSize)).Find(&users)
		}
	} else {
		// show without pagination
		h.logger.Println("* without pagination")
		if filters != nil {
			h.logger.Println("* with filtering")
			h.db.DB.Where(&filters).Find(&users)
		} else {
			h.logger.Println("* without filtering")
			h.db.DB.Find(&users)
		}
	}

	err = json.NewEncoder(w).Encode(&users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("error: users have been listed successfully, but there's a problem returning the response: %v", err)
		return
	}

	h.logger.Println("users have been listed successfully")
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("request list a single user")

	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	id := getIdFromVars(r)
	var user store.User
	h.db.DB.Where("id = ?", id).Find(&user)

	err := json.NewEncoder(w).Encode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("error: user has been listed successfully, but there's a problem returning the response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Println("users have been listed successfully")
}

func (h *Handlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("update user request")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	id := getIdFromVars(r)
	var user store.User
	h.db.DB.Where("id = ?", id).Find(&user)

	if &user == nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Cannot find user with id %s", id)
		h.logger.Println(fmt.Sprintf("error updating user: %s", msg))

		_, err := w.Write([]byte(msg))
		if err != nil {
			h.logger.Println(fmt.Sprintf("error returning the reponse: %v", err))
		}
		return
	}

	newUser, rerr, statusCode := unmarshalUserFromRequestBody(r)
	if rerr != nil {
		w.WriteHeader(statusCode)
		h.logger.Println(fmt.Sprintf("error update user: %s", rerr))
		return
	}
	mergeUserObjects(&user, newUser)

	h.db.DB.Save(user)

	msg := fmt.Sprintf("User with id %s has been updated successfully", id)
	h.logger.Println(msg)

	_, err := w.Write([]byte(msg))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

	}
	w.WriteHeader(http.StatusOK)

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

type ResponseError struct {
	message string
}
func (e *ResponseError) Error() string {
	return e.message
}

type EmptyBody ResponseError

func (e *EmptyBody) Error() string {
	return e.message
}

func unmarshalUserFromRequestBody(r *http.Request) (*store.User, error, int) {
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
			return nil, &ResponseError{fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)}, http.StatusBadRequest

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, &ResponseError{fmt.Sprintf("Request body contains badly-formed JSON") }, http.StatusBadRequest

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			return nil, &ResponseError{fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset) }, http.StatusBadRequest

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return nil, &ResponseError{fmt.Sprintf("Request body contains unknown field %s", fieldName) }, http.StatusBadRequest

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			return nil, &EmptyBody{"Request body must not be empty" }, http.StatusBadRequest

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return nil, &ResponseError{fmt.Sprintf("Request body must not be larger than 1MB") }, http.StatusRequestEntityTooLarge

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return nil, &ResponseError{fmt.Sprintf("Internal server error") }, http.StatusInternalServerError
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return nil,  &ResponseError{fmt.Sprintf("Request body must only contain a single JSON object") }, http.StatusBadRequest
	}

	return &user, nil, http.StatusOK
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
	h.logger.Println("setting router and handle function")
	router := mux.NewRouter()
	router.HandleFunc("/users", h.GetUsers).Methods(http.MethodGet) // without pagination
	router.HandleFunc("/users/{pagination-size:[0-9]+}/{page:[0-9]+}", h.GetUsers).Methods(http.MethodGet) // with pagination
	router.HandleFunc("/user", h.AddUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.UpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", h.DeleteUser).Methods(http.MethodDelete)

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
