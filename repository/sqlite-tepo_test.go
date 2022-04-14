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
	dbName     = "test"
	testRepo   UserRepository
	testLogger = log.New(os.Stdout, "testing-repository", log.LstdFlags|log.Llongfile)
	testUsers  = []model.User{
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
	testRepo, err = NewSqliteRepo(dbName, testLogger)
	require.NoError(t, err)
}

func cleanTestCase(t *testing.T) {
	require.NoError(t, os.Remove(fmt.Sprintf("%s.db", dbName)))
}

func TestNewSqliteRepo(t *testing.T) {
	r, err := NewSqliteRepo(dbName, testLogger)
	require.NoError(t, err)
	require.NotEmpty(t, r)
	require.NoError(t, os.Remove(fmt.Sprintf("%s.db", dbName)))
}

func TestAddOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	result, err := testRepo.Add(&testUsers[0])
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, &testUsers[0], result)

}

func TestDeleteOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testRepo.Add(&testUsers[0])

	err := testRepo.Delete(testUsers[0].ID)
	require.NoError(t, err)
}

func TestGetOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	user, err := testRepo.Add(&testUsers[0])
	require.NoError(t, err)
	require.Equal(t, &testUsers[0], user)
}

func TestGetAllOK(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testRepo.Add(&testUsers[0])
	_, _ = testRepo.Add(&testUsers[1])
	users, err := testRepo.GetAll(nil, 0, 0)
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

func TestUpdate(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	_, _ = testRepo.Add(&testUsers[0])
	user, err := testRepo.Update(&testUsers[0], &testUsers[1])
	require.NoError(t, err)
	require.Equal(t, testUsers[0].ID, user.ID)
	require.Equal(t, testUsers[0].FirstName, user.FirstName)
	require.Equal(t, testUsers[0].LastName, user.LastName)
	require.Equal(t, testUsers[0].Nickname, user.Nickname)
	require.Equal(t, testUsers[0].Password, user.Password)
	require.Equal(t, testUsers[0].Email, user.Email)
	require.Equal(t, testUsers[0].Country, user.Country)
}
