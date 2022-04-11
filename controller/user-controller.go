package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pavelerokhin/cleanarchitecture-restapi-go/errors"
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
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "error unmarshalling the request"})
		return
	}

	errV := c.Service.Validate(&user)
	if errV != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: errV.Error()})
		return
	}

	postC, errC := c.Service.Add(&user)
	if errC != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "error saving user"})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(postC)
}

func (c controller) DeleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	id, err := c.Service.Delete(request)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("{\"error\": \"error while deleting a User with id %v: %v\"}", id, err)
		c.Logger.Println(msg)
		_, _ = response.Write([]byte(msg))

		return
	}

	_, err = response.Write([]byte(fmt.Sprintf("{\"message\": \"user with id %v has beeen deleted successfully\"}", id)))
	if err != nil {
		c.Logger.Printf("error: User has been successfully deleted, but there's a problem returning the response: %v", err)

		response.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("{\"error\": \"User with id %v has been successfully deleted, but there's a problem returning the response: %v\"}", id, err)
		c.Logger.Println(msg)
		_, _ = response.Write([]byte(msg))

		return
	}

	response.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("{\"error\": \"User with id %v has been deleted successfully\"}", id)
	c.Logger.Println(msg)
	_, _ = response.Write([]byte(msg))
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
		json.NewEncoder(response).Encode(errors.ServiceError{Message: fmt.Sprintf("error getting users from the database: %v", err)})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(users)
}

func (c controller) UpdateUser(response http.ResponseWriter, request *http.Request) {
	c.Logger.Println("update user request")
	response.Header().Set("Content-Type", "application/json")

	user, err := c.Service.Update(request)

	if err != nil {
		_, err := response.Write([]byte(fmt.Sprintf("{\"error\":\"error returning the reponse: %v\"}", err)))
		if err != nil {
			c.Logger.Println(fmt.Sprintf("error returning the reponse: %v", err))
			return
		}
	}

	err = json.NewEncoder(response).Encode(user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		c.Logger.Println(fmt.Sprintf("user has been updated, but we had an error returning the reponse to the client: %v", err))
		return
	}
	response.WriteHeader(http.StatusOK)

}
