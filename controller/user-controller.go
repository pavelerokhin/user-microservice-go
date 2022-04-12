package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pavelerokhin/user-microservice-go/model"
	"github.com/pavelerokhin/user-microservice-go/service"
)

type controller struct {
	Logger  *log.Logger
	Service service.UserService
}

type UserController interface {
	AddUser(response http.ResponseWriter, request *http.Request)
	DeleteUser(response http.ResponseWriter, request *http.Request)
	CreateUser(response http.ResponseWriter, request *http.Request)
	GetUser(response http.ResponseWriter, request *http.Request)
	GetAllUsers(response http.ResponseWriter, request *http.Request)
	UpdateUser(response http.ResponseWriter, request *http.Request)
}

func New(service service.UserService, logger *log.Logger) UserController {
	return &controller{Logger: logger, Service: service}
}

func (c controller) AddUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		msg := fmt.Sprintf("error unmarshalling the request: %v", err)
		tryToResponseJsonError(response, c.Logger, msg)
		return
	}

	errValidation := c.Service.Validate(&user)
	if errValidation != nil {
		msg := fmt.Sprintf("error validating the request: %v", errValidation.Error())
		tryToResponseJsonError(response, c.Logger, msg)
		return
	}

	userAdded, errC := c.Service.Add(&user)
	if errC != nil {
		msg := "error saving user"
		tryToResponseJsonError(response, c.Logger, msg)
		return
	}

	tryToResponseUserOK(response, c.Logger, userAdded)
}

func (c controller) DeleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	id, err := c.Service.Delete(request)
	if err != nil {
		msg := fmt.Sprintf("error while deleting a User with id %v: %v", id, err)
		c.Logger.Println(msg)
		tryToResponseJsonError(response, c.Logger, msg)
		return
	}

	msg := fmt.Sprintf("user with id %v has beeen deleted successfully", id)
	tryToResponseMsgOK(response, c.Logger, msg)
}

func (c controller) CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

}

func (c controller) GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

}

func (c controller) GetAllUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	users, err, statusCode := c.Service.GetAll(request)
	if err != nil {
		response.WriteHeader(statusCode)
		msg := fmt.Sprintf("error getting users from the database: %v", err)
		tryToResponseJsonError(response, c.Logger, msg)
		return
	}

	tryToResponseUsersOK(response, c.Logger, users)
}

func (c controller) UpdateUser(response http.ResponseWriter, request *http.Request) {
	c.Logger.Println("update user request")
	response.Header().Set("Content-Type", "application/json")

	user, err, statusCode := c.Service.Update(request)

	if err != nil {
		response.WriteHeader(statusCode)
		msg := fmt.Sprintf("error returning the reponse: %v", err)
		tryToResponseMsgOK(response, c.Logger, msg)
	}

	tryToResponseUserOK(response, c.Logger, user)
}
