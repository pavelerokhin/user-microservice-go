package service

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pavelerokhin/user-microservice-go/model"
)

var (
	logger         = log.New(os.Stdout, "testing-user-service", log.LstdFlags|log.Llongfile)
	mockRepository = new(MockRepository)
	testService    = New(mockRepository, logger)

	users = []model.User{
		{
			ID:        1,
			FirstName: "user1",
			LastName:  "y",
			Nickname:  "z",
			Password:  "1",
			Email:     "a@b.com",
			Country:   "Y",
		},
		{
			ID:        2,
			FirstName: "user2",
			LastName:  "y",
			Nickname:  "z",
			Password:  "1",
			Email:     "a@b.com",
			Country:   "Y",
		},
	}
)

type MockRepository struct {
	mock mock.Mock
}

func (mr *MockRepository) Add(_ *model.User) (*model.User, error) {
	args := mr.mock.Called()
	result := args.Get(0)
	return result.(*model.User), args.Error(1)
}

func (mr *MockRepository) Delete(_ int) error {
	args := mr.mock.Called()
	return args.Error(1)
}

func (mr *MockRepository) Get(_ int) (*model.User, error) {
	args := mr.mock.Called()
	result := args.Get(0)
	return result.(*model.User), args.Error(1)
}

func (mr *MockRepository) GetAll(_ *model.User, _, _ int) ([]model.User, error) {
	args := mr.mock.Called()
	result := args.Get(0)
	return result.([]model.User), args.Error(1)
}

func (mr *MockRepository) Update(_, _ *model.User) (*model.User, error) {
	args := mr.mock.Called()
	result := args.Get(0)
	return result.(*model.User), args.Error(1)
}

// Add function
func TestAdd(t *testing.T) {
	mockRepository.mock.On("Add").Return(&users[0], nil)
	result, err := testService.Add(&users[0])
	// Mock assertion
	mockRepository.mock.AssertExpectations(t)
	// Data assertion
	assert.Equal(t, &users[0], result)
	assert.Nil(t, err)
}

// Delete function
func TestDelete(t *testing.T) {
	mockRepository.mock.On("Delete").Return(1, nil)

	request, _ := http.NewRequest(http.MethodDelete, "/user/1", nil)
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "1",
	}
	request = mux.SetURLVars(request, vars)

	id, err := testService.Delete(request)
	// Mock assertion
	mockRepository.mock.AssertExpectations(t)
	// Data assertion
	assert.Equal(t, 1, id)
	assert.Nil(t, err)
}

// Get function
func TestGet(t *testing.T) {
	mockRepository.mock.On("Get").Return(&users[0], nil)
	request, _ := http.NewRequest(http.MethodGet, "/user/1", nil)
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "1",
	}
	request = mux.SetURLVars(request, vars)

	result, statusCode, err := testService.Get(request)
	// Mock assertion
	mockRepository.mock.AssertExpectations(t)
	// Data assertion
	assert.Equal(t, &users[0], result)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)
}

// GetAll function
func TestGetAll(t *testing.T) {
	mockRepository.mock.On("GetAll").Return(users, nil)
	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	result, statusCode, err := testService.GetAll(request)
	// Mock assertion
	mockRepository.mock.AssertExpectations(t)
	// Data assertion
	assert.Equal(t, 2, len(result))
	assert.Equal(t, users, result)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)
}

// Update function
func TestUpdate(t *testing.T) {
	oldUser := users[0]
	newUser := users[0]
	newUser.FirstName = "updated first name"

	mockRepository.mock.On("Get").Return(&users[0], nil)
	mockRepository.mock.On("Update").Return(&newUser, nil)
	requestBody, err := json.Marshal(map[string]string{"first_name": "updated first name"})
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodGet, "/user/1", bytes.NewBuffer(requestBody))
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "1",
	}
	request = mux.SetURLVars(request, vars)

	//Hack to try to fake gorilla/mux vars
	result, statusCode, err := testService.Update(request)
	// Mock assertion
	mockRepository.mock.AssertExpectations(t)
	// Data assertion
	assert.Equal(t, oldUser.ID, result.ID)
	assert.Equal(t, newUser.FirstName, result.FirstName)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)
}

// Validate method
func TestValidateEmptyUser(t *testing.T) {
	err := testService.Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "the user object is empty", err.Error())
}

func TestValidateEmptyUserFirstName(t *testing.T) {
	user := model.User{
		FirstName: "",
	}
	err := testService.Validate(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "the user's first name is empty", err.Error())
}

func TestValidateEmptyUserLastName(t *testing.T) {
	user := model.User{
		FirstName: "x",
		LastName:  "",
	}
	err := testService.Validate(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "the user's last name is empty", err.Error())
}

func TestValidateEmptyUserNickname(t *testing.T) {
	user := model.User{
		FirstName: "x",
		LastName:  "y",
		Nickname:  "",
	}
	err := testService.Validate(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "the user's nickname is empty", err.Error())
}

func TestValidateEmptyUserPassword(t *testing.T) {
	user := model.User{
		FirstName: "x",
		LastName:  "y",
		Nickname:  "z",
		Password:  "",
	}
	err := testService.Validate(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "the user's password is empty", err.Error())
}

func TestValidateEmptyUserEmail(t *testing.T) {
	user := model.User{
		FirstName: "x",
		LastName:  "y",
		Nickname:  "z",
		Password:  "1",
		Email:     "",
	}
	err := testService.Validate(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "the user's email field is empty", err.Error())
}

func TestValidateEmptyUserCountry(t *testing.T) {
	user := model.User{
		FirstName: "x",
		LastName:  "y",
		Nickname:  "z",
		Password:  "1",
		Email:     "a@b.com",
		Country:   "",
	}
	err := testService.Validate(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "the user's country field is empty", err.Error())
}

func TestValidateUserOK(t *testing.T) {
	user := model.User{
		FirstName: "x",
		LastName:  "y",
		Nickname:  "z",
		Password:  "1",
		Email:     "a@b.com",
		Country:   "Y",
	}
	err := testService.Validate(&user)
	assert.Nil(t, err)
}
