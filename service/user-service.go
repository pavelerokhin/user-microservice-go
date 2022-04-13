// pkg represents the service layer of the microservice

package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/pavelerokhin/user-microservice-go/errs"
	"github.com/pavelerokhin/user-microservice-go/model"
	"github.com/pavelerokhin/user-microservice-go/repository"
)

type UserService interface {
	Add(user *model.User) (*model.User, error)
	Delete(request *http.Request) (int, error)
	Get(request *http.Request) (*model.User, error, int)
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

func (s *service) Get(request *http.Request) (*model.User, error, int) {
	s.Logger.Println("request get single user")

	id, err := getIdFromRequestVars(request)
	if err != nil {
		return nil, fmt.Errorf("error while parsing user's ID: %v", err), http.StatusInternalServerError
	}

	user, err := s.Repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving user with ID %v: %v", id, err), http.StatusInternalServerError
	}

	s.Logger.Println(fmt.Sprintf("user with ID %v has been retrieved successfully", id))
	return user, err, http.StatusOK
}

func (s *service) GetAll(request *http.Request) ([]model.User, error, int) {
	s.Logger.Println("request list users")
	var filters *model.User

	filters, err, statusCode := unmarshalUserFromRequest(request)

	errEmptyBody := &errs.EmptyBody{}
	if err != nil && !errors.As(err, &errEmptyBody) {
		return nil, fmt.Errorf("error while parsing filter parameters: %v", err), statusCode
	}

	var pageSize, page int
	vars := mux.Vars(request)
	if vars["page-size"] != "" {
		pageSize, err = strconv.Atoi(vars["page-size"])
		if err != nil {
			msg := fmt.Errorf("cannot get pagination limit: %v", err)
			return nil, msg, http.StatusInternalServerError
		}

		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			msg := fmt.Errorf("cannot get page: %v", err)
			return nil, msg, http.StatusInternalServerError
		}

		if pageSize <= 0 {
			msg := fmt.Errorf("page size cannot be less then 1")
			return nil, msg, http.StatusBadRequest
		}

		if page <= 0 {
			msg := fmt.Errorf("cannot get page less then 1")
			return nil, msg, http.StatusBadRequest
		}
	}

	msg := "try to list users "
	if pageSize == 0 {
		msg += "without pagination "
	} else {
		msg += "with pagination "
	}
	if filters == nil {
		msg += "without filtering"
	} else {
		msg += "with filtering"
	}

	s.Logger.Printf(msg)

	allUsers, err := s.Repo.GetAll(filters, pageSize, page)
	if err != nil {
		return allUsers, err, http.StatusInternalServerError
	}
	return allUsers, err, http.StatusOK
}

func (s *service) Update(request *http.Request) (*model.User, error, int) {
	s.Logger.Println("request update a user")

	// parse user
	id, err := getIdFromRequestVars(request)
	var user *model.User

	user, err = s.Repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error while trying to find the user to update (ID %v): %v", id, err),
			http.StatusBadRequest
	}

	if user == nil {
		return nil, fmt.Errorf("error updating user: cannot find user with ID %v", id),
			http.StatusBadRequest
	}

	newUser, err, statusCode := unmarshalUserFromRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %s", err), statusCode
	}

	//mergeUserObjects(user, newUser)

	user, err = s.Repo.Update(user, newUser)
	if err != nil {
		return nil, fmt.Errorf("error while updating user with ID %v: %v", id, err),
			http.StatusInternalServerError
	}

	s.Logger.Println(fmt.Errorf("user with ID %v has been updated successfully", id))
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
