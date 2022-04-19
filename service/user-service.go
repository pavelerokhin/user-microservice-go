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
	Get(request *http.Request) (*model.User, int, error)
	GetAll(request *http.Request) ([]model.User, int, error)
	Update(request *http.Request) (*model.User, int, error)
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
	s.Logger.Println("service request add a new user")
	return s.Repo.Add(user)
}

func (s *service) Delete(request *http.Request) (int, error) {
	s.Logger.Println("service request delete user")

	id, err := getIDFromRequestVars(request)
	if id == 0 || err != nil {
		return 0, fmt.Errorf("cannot parse ID of the user to delete: %w", err)
	}

	return id, s.Repo.Delete(id)
}

func (s *service) Get(request *http.Request) (*model.User, int, error) {
	s.Logger.Println("service request get single user")

	id, err := getIDFromRequestVars(request)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error while parsing user's ID: %v", err)
	}

	user, err := s.Repo.Get(id)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error while retrieving user with ID %v: %v", id, err)
	}

	s.Logger.Println(fmt.Sprintf("user with ID %v has been retrieved successfully", id))
	return user, http.StatusOK, err
}

func (s *service) GetAll(request *http.Request) ([]model.User, int, error) {
	s.Logger.Println("service request list users")
	var filters *model.User

	filters, err, statusCode := unmarshalUserFromRequest(request)

	errEmptyBody := &errs.EmptyBody{}
	if err != nil && !errors.As(err, &errEmptyBody) {
		return nil, statusCode, fmt.Errorf("error while parsing filter parameters: %v", err)
	}

	var pageSize, page int
	vars := mux.Vars(request)
	if vars["page-size"] != "" {
		pageSize, err = strconv.Atoi(vars["page-size"])
		if err != nil {
			msg := fmt.Errorf("cannot get pagination limit: %v", err)
			return nil, http.StatusInternalServerError, msg
		}

		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			msg := fmt.Errorf("cannot get page: %v", err)
			return nil, http.StatusInternalServerError, msg
		}

		if pageSize <= 0 {
			msg := fmt.Errorf("page size cannot be less then 1")
			return nil, http.StatusBadRequest, msg
		}

		if page <= 0 {
			msg := fmt.Errorf("cannot get page less then 1")
			return nil, http.StatusBadRequest, msg
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
		return allUsers, http.StatusInternalServerError, err
	}
	return allUsers, http.StatusOK, err
}

func (s *service) Update(request *http.Request) (*model.User, int, error) {
	s.Logger.Println("service request update a user")

	// parse user
	var id, err = getIDFromRequestVars(request)
	var user *model.User

	user, err = s.Repo.Get(id)
	if err != nil {
		return nil,
			http.StatusBadRequest,
			fmt.Errorf("error while trying to find the user to update (ID %v): %v", id, err)
	}

	if user == nil {
		return nil,
			http.StatusBadRequest,
			fmt.Errorf("error updating user: cannot find user with ID %v", id)
	}

	newUser, err, statusCode := unmarshalUserFromRequest(request)
	if err != nil {
		return nil, statusCode, fmt.Errorf("error updating user: %s", err)
	}

	user, err = s.Repo.Update(user, newUser)
	if err != nil {
		return nil,
			http.StatusInternalServerError,
			fmt.Errorf("error while updating user with ID %v: %v", id, err)
	}

	s.Logger.Println(fmt.Errorf("user with ID %v has been updated successfully", id))
	return user, http.StatusOK, nil
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

	if !user.CreatedAt.IsZero() {
		err := errors.New("the user's create time must be empty")
		return err
	}
	if !user.UpdatedAt.IsZero() {
		err := errors.New("the user's update time must be empty")
		return err
	}

	return nil
}
