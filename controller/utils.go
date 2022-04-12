package controller

import (
	"encoding/json"
	"github.com/pavelerokhin/user-microservice-go/errs"
	"github.com/pavelerokhin/user-microservice-go/model"
	"log"
	"net/http"
)

const (
	errMsgEncodeOK = "error while encoding the response from the server (the user request has been processed)"
	errMsgEncodeKO = "error while encoding the response from the server (the user request hasn't been processed)"
)

func writeResponseJSON(response http.ResponseWriter, msg string) error {
	return json.NewEncoder(response).Encode(errs.ResponseError{Message: msg})
}

func tryToResponseJSONError(response http.ResponseWriter, logger *log.Logger, msg string, statusCode int) {
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

// this function was duplicated to avoid using reflection
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

// this function was duplicated to avoid using reflection
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

// this function was duplicated to avoid using reflection
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
