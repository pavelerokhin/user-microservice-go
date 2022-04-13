package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pavelerokhin/user-microservice-go/errs"
	"github.com/pavelerokhin/user-microservice-go/model"
)

const (
	errMsgEncodeOK = "error while encoding the response from the server (the user request has been processed)"
	errMsgEncodeKO = "error while encoding the response from the server (the user request hasn't been processed)"
)

func writeResponseJSON(response http.ResponseWriter, msg string) error {
	return json.NewEncoder(response).Encode(errs.ResponseError{Message: msg})
}

// tryToResponseJSONError is a utility function that tries to write a response (error) message formatted
// as an JSON to the client. In case it couldn't write the message, it tries to return a standard errMsgEncodeKO
// message to the client. statusCode of the response can be specified. In case statusCode is 0,
// the response is considered to be http.StatusInternalServerError (500)
func tryToResponseJSONError(response http.ResponseWriter, logger *log.Logger, statusCode int, msg string) {
	logger.Println(msg)
	if statusCode == 0 {
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.WriteHeader(statusCode)
	}

	err := writeResponseJSON(response, msg)
	if err != nil {
		logger.Println(errMsgEncodeKO)
		_ = writeResponseJSON(response, errMsgEncodeKO)
		return
	}
}

// tryToResponseMsgOK is similar to tryToResponseMsgKO, but is supposed to return the response message
// in cases when no error has occurred. This duplicates tryToResponseMsgOK to avoid using reflection
func tryToResponseMsgOK(response http.ResponseWriter, logger *log.Logger, msg string) {
	logger.Println(msg)
	err := writeResponseJSON(response, msg)
	if err != nil {
		logger.Println(errMsgEncodeOK)
		response.WriteHeader(http.StatusInternalServerError)
		_ = writeResponseJSON(response, errMsgEncodeOK)
		return
	}
	response.WriteHeader(http.StatusOK)
}

// tryToResponseUserOK duplicates tryToResponseMsgOK; it marshals the User object in the response
func tryToResponseUserOK(response http.ResponseWriter, logger *log.Logger, msg *model.User) {
	logger.Println(msg)
	err := json.NewEncoder(response).Encode(msg)
	if err != nil {
		logger.Println(errMsgEncodeOK)
		response.WriteHeader(http.StatusInternalServerError)
		_ = writeResponseJSON(response, errMsgEncodeOK)
		return
	}
	response.WriteHeader(http.StatusOK)
}

// tryToResponseUserOK duplicates tryToResponseMsgOK; it marshals the slice of User objects in the response
func tryToResponseUsersOK(response http.ResponseWriter, logger *log.Logger, msg []model.User) {
	logger.Println(msg)
	err := json.NewEncoder(response).Encode(msg)
	if err != nil {
		logger.Println(errMsgEncodeOK)
		response.WriteHeader(http.StatusInternalServerError)
		_ = writeResponseJSON(response, errMsgEncodeOK)
		return
	}
	response.WriteHeader(http.StatusOK)
}
