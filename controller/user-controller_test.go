package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/pavelerokhin/user-microservice-go/model"
	"github.com/pavelerokhin/user-microservice-go/repository"
	"github.com/pavelerokhin/user-microservice-go/service"
)

var (
	repositoryName = "user-controller-testing"

	testLogger = log.New(os.Stdout, "testing-controller", log.LstdFlags|log.Llongfile)
	testUser   = model.User{
		ID:        1,
		FirstName: "user1",
		LastName:  "y",
		Nickname:  "z",
		Password:  "1",
		Email:     "a@b.com",
		Country:   "Y",
	}

	testUserRepository repository.UserRepository
	testUserService    service.UserService
	testUserController UserController
)

func setupTestCase(t *testing.T) {
	var err error
	testUserRepository, err = repository.NewSqliteRepo(repositoryName, testLogger)
	testUserService = service.New(testUserRepository, testLogger)
	testUserController = New(testUserService, testLogger)
	require.NoError(t, err)
}

func setupTestCaseWithUser(t *testing.T) {
	var err error
	testUserRepository, err = repository.NewSqliteRepo(repositoryName, testLogger)
	testUserService = service.New(testUserRepository, testLogger)
	testUserController = New(testUserService, testLogger)
	require.NoError(t, err)
	_, err = testUserRepository.Add(&testUser)
	require.NoError(t, err)
}

func cleanTestCase(t *testing.T) {
	require.NoError(t, os.Remove(fmt.Sprintf("%s.db", repositoryName)))
}

func TestMain(m *testing.M) {
	_ = os.Remove(fmt.Sprintf("%s.db", repositoryName))
	code := m.Run()
	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	// Create a new HTTP POST request
	jsonUser, err := json.Marshal(testUser)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}

	// Assign HTTP Handle function (controller AddUser function)
	handler := http.HandlerFunc(testUserController.AddUser)

	// Record HTTP Response (httptest library)
	response := httptest.NewRecorder()

	// Dispatch the HTTP request
	handler.ServeHTTP(response, request)

	// Add assertions on the HTTP status code and the response
	status := response.Code
	require.Equal(t, http.StatusOK, status)

	// Decode HTTP response
	var user model.User
	err = json.NewDecoder(io.Reader(response.Body)).Decode(&user)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, testUser.FirstName, user.FirstName)
	require.Equal(t, testUser.LastName, user.LastName)
	require.Equal(t, testUser.Nickname, user.Nickname)
	require.Equal(t, testUser.Password, user.Password)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.Country, user.Country)
}

func TestDeleteUser(t *testing.T) {
	setupTestCaseWithUser(t)
	defer cleanTestCase(t)

	// Create a new HTTP POST request to delete user
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", testUser.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": strconv.Itoa(testUser.ID),
	}
	request = mux.SetURLVars(request, vars)

	// Assign HTTP Handle function (controller AddUser function)
	handler := http.HandlerFunc(testUserController.DeleteUser)

	// Record HTTP Response (httptest library)
	response := httptest.NewRecorder()

	// Dispatch the HTTP request
	handler.ServeHTTP(response, request)

	// Add assertions on the HTTP status code and the response
	status := response.Code
	require.Equal(t, http.StatusOK, status)

	// Try to get user id 1 (should be nil)
	user, err := testUserRepository.Get(testUser.ID)
	require.Error(t, err)
	require.Nil(t, user)
}

func TestGetUser(t *testing.T) {
	setupTestCaseWithUser(t)
	defer cleanTestCase(t)

	// Create a new HTTP POST request to delete user
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", testUser.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": strconv.Itoa(testUser.ID),
	}
	request = mux.SetURLVars(request, vars)

	// Assign HTTP Handle function (controller AddUser function)
	handler := http.HandlerFunc(testUserController.GetUser)

	// Record HTTP Response (httptest library)
	response := httptest.NewRecorder()

	// Dispatch the HTTP request
	handler.ServeHTTP(response, request)

	// Add assertions on the HTTP status code and the response
	status := response.Code
	require.Equal(t, http.StatusOK, status)

	// Decode HTTP response
	var user model.User
	err = json.NewDecoder(io.Reader(response.Body)).Decode(&user)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, testUser.FirstName, user.FirstName)
	require.Equal(t, testUser.LastName, user.LastName)
	require.Equal(t, testUser.Nickname, user.Nickname)
	require.Equal(t, testUser.Password, user.Password)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.Country, user.Country)
}

func TestGetAllUsers(t *testing.T) {
	setupTestCaseWithUser(t)
	defer cleanTestCase(t)

	// Create a new HTTP POST request to delete user
	request, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Assign HTTP Handle function (controller AddUser function)
	handler := http.HandlerFunc(testUserController.GetAllUsers)

	// Record HTTP Response (httptest library)
	response := httptest.NewRecorder()

	// Dispatch the HTTP request
	handler.ServeHTTP(response, request)

	// Add assertions on the HTTP status code and the response
	status := response.Code
	require.Equal(t, http.StatusOK, status)

	// Decode HTTP response
	var users []model.User
	err = json.NewDecoder(io.Reader(response.Body)).Decode(&users)
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, testUser.FirstName, users[0].FirstName)
	require.Equal(t, testUser.LastName, users[0].LastName)
	require.Equal(t, testUser.Nickname, users[0].Nickname)
	require.Equal(t, testUser.Password, users[0].Password)
	require.Equal(t, testUser.Email, users[0].Email)
	require.Equal(t, testUser.Country, users[0].Country)
}

func TestUpdateUser(t *testing.T) {
	setupTestCaseWithUser(t)
	defer cleanTestCase(t)

	// Create a new HTTP POST request to delete user

	requestBody, err := json.Marshal(map[string]string{"first_name": "updated first name"})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/%d", testUser.ID), bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	vars := map[string]string{
		"id": strconv.Itoa(testUser.ID),
	}
	request = mux.SetURLVars(request, vars)

	// Assign HTTP Handle function (controller AddUser function)
	handler := http.HandlerFunc(testUserController.UpdateUser)

	// Record HTTP Response (httptest library)
	response := httptest.NewRecorder()

	// Dispatch the HTTP request
	handler.ServeHTTP(response, request)

	// Add assertions on the HTTP status code and the response
	status := response.Code
	require.Equal(t, http.StatusOK, status)

	// Decode HTTP response
	var user model.User
	err = json.NewDecoder(io.Reader(response.Body)).Decode(&user)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "updated first name", user.FirstName)
	require.Equal(t, testUser.LastName, user.LastName)
	require.Equal(t, testUser.Nickname, user.Nickname)
	require.Equal(t, testUser.Password, user.Password)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.Country, user.Country)
}
