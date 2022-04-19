package repository

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavelerokhin/user-microservice-go/model"
)

var (
	dbName             = "test"
	testUserRepository UserRepository
	testLogger         = log.New(os.Stdout, "testing-repository", log.LstdFlags|log.Llongfile)
	testUsers          = []model.User{
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

func setupTestCase(t *testing.T) {
	var err error
	testUserRepository, err = NewSqliteRepo(dbName, testLogger)
	require.NoError(t, err)
}

func cleanTestCase(t *testing.T) {
	require.NoError(t, os.Remove(fmt.Sprintf("%s.db", dbName)))
}

// NewSqliteRepo function testing
func TestNewSqliteRepoOK(t *testing.T) {
	r, err := NewSqliteRepo(dbName, testLogger)
	require.NoError(t, err)
	require.NotEmpty(t, r)
	require.NoError(t, os.Remove(fmt.Sprintf("%s.db", dbName)))
}

func TestNewSqliteRepoEmptyNameKO(t *testing.T) {
	r, err := NewSqliteRepo("", testLogger)
	require.Error(t, err)
	require.Equal(t, "database name is empty", err.Error())
	require.Empty(t, r)
}

// Add function testing
func TestAddOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)
	result, err := testUserRepository.Add(&testUsers[0])
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, &testUsers[0], result)
}

func TestAddNotUniqueKO(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)
	result, err := testUserRepository.Add(&testUsers[0])
	require.NoError(t, err)
	require.NotEmpty(t, result)

	result, err = testUserRepository.Add(&testUsers[0])
	require.Error(t, err)
	require.Nil(t, result)
}

// Delete function testing
func TestDeleteOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testUserRepository.Add(&testUsers[0])

	err := testUserRepository.Delete(testUsers[0].ID)
	require.NoError(t, err)
}

func TestDeleteNoIdKO(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	err := testUserRepository.Delete(testUsers[0].ID)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("error: cannot find user with ID %v", testUsers[0].ID),
		err.Error())
}

// Get function testing
func TestGetOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	user, err := testUserRepository.Add(&testUsers[0])
	require.NoError(t, err)
	require.Equal(t, &testUsers[0], user)

	userGet, errGet := testUserRepository.Get(testUsers[0].ID)
	require.NoError(t, errGet)
	require.Equal(t, user.ID, userGet.ID)
	require.Equal(t, user.FirstName, userGet.FirstName)
	require.Equal(t, user.LastName, userGet.LastName)
	require.Equal(t, user.Nickname, userGet.Nickname)
	require.Equal(t, user.Password, userGet.Password)
	require.Equal(t, user.Email, userGet.Email)
}

func TestGetNoIdKO(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	user, err := testUserRepository.Get(testUsers[0].ID)
	require.Nil(t, user)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("user with ID %v not found", testUsers[0].ID), err.Error())
}

// GetAll function testing
func TestGetAllOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testUserRepository.Add(&testUsers[0])
	_, _ = testUserRepository.Add(&testUsers[1])
	users, err := testUserRepository.GetAll(nil, 0, 0)
	require.NoError(t, err)
	require.Equal(t, testUsers[0].ID, users[0].ID)
	require.Equal(t, testUsers[0].FirstName, users[0].FirstName)
	require.Equal(t, testUsers[0].LastName, users[0].LastName)
	require.Equal(t, testUsers[0].Nickname, users[0].Nickname)
	require.Equal(t, testUsers[0].Password, users[0].Password)
	require.Equal(t, testUsers[0].Email, users[0].Email)
	require.Equal(t, testUsers[0].Country, users[0].Country)
	require.Equal(t, testUsers[1].ID, users[1].ID)
	require.Equal(t, testUsers[1].FirstName, users[1].FirstName)
	require.Equal(t, testUsers[1].LastName, users[1].LastName)
	require.Equal(t, testUsers[1].Nickname, users[1].Nickname)
	require.Equal(t, testUsers[1].Password, users[1].Password)
	require.Equal(t, testUsers[1].Email, users[1].Email)
	require.Equal(t, testUsers[1].Country, users[1].Country)
}

func TestGetAllPaginationOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testUserRepository.Add(&testUsers[0])
	_, _ = testUserRepository.Add(&testUsers[1])
	users, err := testUserRepository.GetAll(nil, 1, 1)
	require.NoError(t, err)
	require.Equal(t, testUsers[0].ID, users[0].ID)
	require.Equal(t, testUsers[0].FirstName, users[0].FirstName)
	require.Equal(t, testUsers[0].LastName, users[0].LastName)
	require.Equal(t, testUsers[0].Nickname, users[0].Nickname)
	require.Equal(t, testUsers[0].Password, users[0].Password)
	require.Equal(t, testUsers[0].Email, users[0].Email)
	require.Equal(t, testUsers[0].Country, users[0].Country)

	users, err = testUserRepository.GetAll(nil, 1, 2)
	require.Equal(t, testUsers[1].ID, users[0].ID)
	require.Equal(t, testUsers[1].FirstName, users[0].FirstName)
	require.Equal(t, testUsers[1].LastName, users[0].LastName)
	require.Equal(t, testUsers[1].Nickname, users[0].Nickname)
	require.Equal(t, testUsers[1].Password, users[0].Password)
	require.Equal(t, testUsers[1].Email, users[0].Email)
	require.Equal(t, testUsers[1].Country, users[0].Country)
}

func TestGetAllFilteringOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testUserRepository.Add(&testUsers[0])
	_, _ = testUserRepository.Add(&testUsers[1])
	users, err := testUserRepository.GetAll(&testUsers[0], 0, 0)
	require.NoError(t, err)
	require.Equal(t, testUsers[0].ID, users[0].ID)
	require.Equal(t, testUsers[0].FirstName, users[0].FirstName)
	require.Equal(t, testUsers[0].LastName, users[0].LastName)
	require.Equal(t, testUsers[0].Nickname, users[0].Nickname)
	require.Equal(t, testUsers[0].Password, users[0].Password)
	require.Equal(t, testUsers[0].Email, users[0].Email)
	require.Equal(t, testUsers[0].Country, users[0].Country)
}

func TestGetAllPaginationFilteringOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testUserRepository.Add(&testUsers[0])
	_, _ = testUserRepository.Add(&testUsers[1])
	users, err := testUserRepository.GetAll(&model.User{
		ID: testUsers[1].ID,
	}, 1, 1)
	require.NoError(t, err)
	require.Equal(t, testUsers[1].ID, users[0].ID)
	require.Equal(t, testUsers[1].FirstName, users[0].FirstName)
	require.Equal(t, testUsers[1].LastName, users[0].LastName)
	require.Equal(t, testUsers[1].Nickname, users[0].Nickname)
	require.Equal(t, testUsers[1].Password, users[0].Password)
	require.Equal(t, testUsers[1].Email, users[0].Email)
	require.Equal(t, testUsers[1].Country, users[0].Country)
}

// Update function testing
func TestUpdateOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testUserRepository.Add(&testUsers[0])
	user, err := testUserRepository.Update(&testUsers[0], &testUsers[1])
	require.NoError(t, err)
	require.Equal(t, testUsers[0].ID, user.ID)
	require.Equal(t, testUsers[0].FirstName, user.FirstName)
	require.Equal(t, testUsers[0].LastName, user.LastName)
	require.Equal(t, testUsers[0].Nickname, user.Nickname)
	require.Equal(t, testUsers[0].Password, user.Password)
	require.Equal(t, testUsers[0].Email, user.Email)
	require.Equal(t, testUsers[0].Country, user.Country)
}
