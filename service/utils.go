package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/pavelerokhin/user-microservice-go/errs"
	"github.com/pavelerokhin/user-microservice-go/model"
)

func getIdFromRequestVars(request *http.Request) (int, error) {
	vars := mux.Vars(request)
	return strconv.Atoi(vars["id"])
}

func mergeUserObjects(userSource, userDest *model.User) {
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

func unmarshalUserFromRequestBody(r *http.Request) (*model.User, error, int) {
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var user model.User
	err := dec.Decode(&user)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errs in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			return nil, &errs.ResponseError{Message: fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)}, http.StatusBadRequest

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errs in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, &errs.ResponseError{Message: fmt.Sprintf("Request body contains badly-formed JSON")}, http.StatusBadRequest

		// Catch any type errs, like trying to assign a string in the
		// JSON request body to an int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			return nil, &errs.ResponseError{Message: fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)}, http.StatusBadRequest

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return nil, &errs.ResponseError{Message: fmt.Sprintf("Request body contains unknown field %s", fieldName)}, http.StatusBadRequest

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			return nil, &errs.EmptyBody{Message: "Request body must not be empty"}, http.StatusBadRequest

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return nil, &errs.ResponseError{Message: fmt.Sprintf("Request body must not be larger than 1MB")}, http.StatusRequestEntityTooLarge

		// Otherwise, default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return nil, &errs.ResponseError{Message: fmt.Sprintf("Internal server error")}, http.StatusInternalServerError
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return nil, &errs.ResponseError{Message: fmt.Sprintf("Request body must only contain a single JSON object")}, http.StatusBadRequest
	}

	return &user, nil, http.StatusOK
}

func withWithoutSuffix(value interface{}) string {
	if value == nil || value == 0 {
		return "without"
	}

	return "with"
}
