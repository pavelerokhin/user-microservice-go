// pkg implements the controller layer

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
		tryToResponseJSONError(response, c.Logger, 0, msg)
		return
	}

	errValidation := c.Service.Validate(&user)
	if errValidation != nil {
		msg := fmt.Sprintf("error validating the request: %v", errValidation.Error())
		tryToResponseJSONError(response, c.Logger, 0, msg)
		return
	}

	userAdded, errC := c.Service.Add(&user)
	if errC != nil {
		msg := "error saving user"
		tryToResponseJSONError(response, c.Logger, 0, msg)
		return
	}

	tryToResponseUserOK(response, c.Logger, userAdded)
}

func (c controller) DeleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	id, err := c.Service.Delete(request)
	if err != nil {
		msg := fmt.Sprintf("error while deleting a User with ID %v: %v", id, err)
		tryToResponseJSONError(response, c.Logger, 0, msg)
		return
	}

	msg := fmt.Sprintf("user with ID %v has beeen deleted successfully", id)
	tryToResponseMsgOK(response, c.Logger, msg)
}

func (c controller) GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user, statusCode, err := c.Service.Get(request)
	if err != nil {
		msg := fmt.Sprintf("error getting user from the database: %v", err)
		tryToResponseJSONError(response, c.Logger, statusCode, msg)
		return
	}

	tryToResponseUserOK(response, c.Logger, user)
}

func (c controller) GetAllUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	users, statusCode, err := c.Service.GetAll(request)
	if err != nil {
		msg := fmt.Sprintf("error getting users from the database: %v", err)
		tryToResponseJSONError(response, c.Logger, statusCode, msg)
		return
	}

	tryToResponseUsersOK(response, c.Logger, users)
}

func (c controller) UpdateUser(response http.ResponseWriter, request *http.Request) {
	c.Logger.Println("update user request")
	response.Header().Set("Content-Type", "application/json")

	user, statusCode, err := c.Service.Update(request)

	if err != nil {
		response.WriteHeader(statusCode)
		msg := fmt.Sprintf("error returning the reponse: %v", err)
		tryToResponseMsgOK(response, c.Logger, msg)
	}

	tryToResponseUserOK(response, c.Logger, user)
}
