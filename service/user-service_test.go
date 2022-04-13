package service

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var logger = log.New(os.Stdout, "testing-user-service", log.LstdFlags|log.Llongfile)

func TestValidateEmptyUser(t *testing.T) {
	testService := New(nil, logger)

	err := testService.Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "the user object is empty", err.Error())
}
