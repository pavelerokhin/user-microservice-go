// pkg represents the repository layer of the microservice.
// SQLite and MongoDB databases are available

package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/pavelerokhin/user-microservice-go/model"
)

type UserRepository interface {
	Add(user *model.User) (*model.User, error)
	Delete(id int) error
	Get(id int) (*model.User, error)
	GetAll(filters *model.User, pageSize, page int) ([]model.User, error)
	Update(user, newUser *model.User) (*model.User, error)
}

type repo struct {
	DB     *gorm.DB
	Logger *log.Logger
}
