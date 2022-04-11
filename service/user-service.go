package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/pavelerokhin/user-microservice-go/errs"
	"github.com/pavelerokhin/user-microservice-go/model"
	"github.com/pavelerokhin/user-microservice-go/repository"
)

var ()

type UserService interface {
	Add(user *model.User) (*model.User, error)
	Delete(request *http.Request) (int, error)
	Get(id int) (*model.User, error, int)
	GetAll(request *http.Request) ([]model.User, error, int)
	Update(request *http.Request) (*model.User, error, int)
	Validate(user *model.User) error
}

type service struct {
	Logger *log.Logger
	Repo   repository.UserRepository
}

func New(repository repository.UserRepository, logger *log.Logger) UserService {
	return &service{Repo: repository, Logger: logger}
}

func (s *service) Add(user *model.User) (*model.User, error) {
	s.Logger.Println("request add a new user")
	user.ID = rand.Int()
	return s.Repo.Add(user)
}

func (s *service) Delete(request *http.Request) (int, error) {
	s.Logger.Println("request delete user")

	id, err := getIdFromRequestVars(request)
	if id == 0 || err != nil {
		return 0, fmt.Errorf("cannot parse ID of the user to delete: %v", err)
	}

	return id, s.Repo.Delete(id)
}

func (s *service) Get(id int) (*model.User, error, int) {
	s.Logger.Println("request get single user")

	user, err := s.Repo.Get(id)

	if err != nil {
		return user, err, http.StatusInternalServerError
	}

	return user, err, http.StatusOK
}

func (s *service) GetAll(request *http.Request) ([]model.User, error, int) {
	s.Logger.Println("request list users")
	var filters *model.User

	filters, rerr, statusCode := unmarshalUserFromRequestBody(request)
	if rerr != nil && !errors.As(rerr, &errs.EmptyBody{}) {
		return nil, fmt.Errorf("error while parsing filter parameters: %v", rerr), statusCode
	}

	vars := mux.Vars(request)

	pageSize, err := strconv.Atoi(vars["page-size"])
	if err != nil {
		msg := fmt.Errorf("cannot get pagination limit: %v", err)
		s.Logger.Println(msg)
		return nil, msg, http.StatusInternalServerError
	}

	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		msg := fmt.Errorf("cannot get page: %v", err)
		s.Logger.Println(msg)
		return nil, msg, http.StatusInternalServerError
	}

	if pageSize <= 0 {
		pageSize = 1
	}

	if page <= 0 {
		pageSize = 1
	}

	s.Logger.Printf("try to list users %v pagination % filtering",
		withWithoutInterfix(pageSize),
		withWithoutInterfix(filters))

	allUsers, err := s.Repo.GetAll(filters, pageSize, page)
	return allUsers, err, http.StatusOK
}

func (s *service) Update(request *http.Request) (*model.User, error, int) {
	s.Logger.Println("request update a user")

	// parse user
	id, err := getIdFromRequestVars(request)
	var user *model.User
	user, err = s.Repo.Get(id)

	if user == nil {
		msg := fmt.Errorf("error updating user: cannot find user with id %v", id)
		s.Logger.Println(msg)

		return nil, msg, http.StatusBadRequest
	}

	newUser, rerr, statusCode := unmarshalUserFromRequestBody(request)
	if rerr != nil {
		msg := fmt.Errorf("error update user: %s", rerr)
		s.Logger.Println(fmt.Sprintf(msg.Error()))

		return nil, msg, statusCode
	}

	mergeUserObjects(user, newUser)

	user, err = s.Repo.Add(user)
	if err != nil {
		return nil, fmt.Errorf("user with id %s has been updated successfully", id), http.StatusInternalServerError

	}

	return user, nil, http.StatusOK
}

func (*service) Validate(user *model.User) error {
	if user == nil {
		err := errors.New("the user object is empty")
		return err
	}
	if user.FirstName == "" {
		err := errors.New("the user's first name is empty")
		return err
	}
	if user.LastName == "" {
		err := errors.New("the user's last name is empty")
		return err
	}
	if user.Nickname == "" {
		err := errors.New("the user's nickname is empty")
		return err
	}
	if user.Password == "" {
		err := errors.New("the user's password is empty")
		return err
	}
	if user.Email == "" {
		err := errors.New("the user's email field is empty")
		return err
	}
	if user.Country == "" {
		err := errors.New("the user's country field is empty")
		return err
	}

	return nil
}

// UTILS
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
			return nil, &errs.ResponseError{fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)}, http.StatusBadRequest

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errs in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, &errs.ResponseError{fmt.Sprintf("Request body contains badly-formed JSON")}, http.StatusBadRequest

		// Catch any type errs, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			return nil, &errs.ResponseError{fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)}, http.StatusBadRequest

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return nil, &errs.ResponseError{fmt.Sprintf("Request body contains unknown field %s", fieldName)}, http.StatusBadRequest

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			return nil, &errs.EmptyBody{"Request body must not be empty"}, http.StatusBadRequest

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return nil, &errs.ResponseError{fmt.Sprintf("Request body must not be larger than 1MB")}, http.StatusRequestEntityTooLarge

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return nil, &errs.ResponseError{fmt.Sprintf("Internal server error")}, http.StatusInternalServerError
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return nil, &errs.ResponseError{fmt.Sprintf("Request body must only contain a single JSON object")}, http.StatusBadRequest
	}

	return &user, nil, http.StatusOK
}

func withWithoutInterfix(value interface{}) string {
	if value == nil || value == 0 {
		return "without"
	}

	return "with"
}
