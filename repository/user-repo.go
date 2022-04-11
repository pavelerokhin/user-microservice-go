package repository

import (
	"github.com/pavelerokhin/user-microservice-go/model"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	Add(user *model.User) (*model.User, error)
	Delete(id int) error
	Get(id int) (*model.User, error)
	GetAll(filters *model.User, pageSize, page int) ([]model.User, error)
	Update(user *model.User) (*model.User, error)
}

type repo struct {
	DB     *gorm.DB
	Logger *log.Logger
}
