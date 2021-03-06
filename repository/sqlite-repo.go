package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/pavelerokhin/user-microservice-go/model"
)

func NewSqliteRepo(dbName string, l *log.Logger) (UserRepository, error) {
	l.Println("preparing SQLite database")

	if dbName == "" {
		return nil, fmt.Errorf("database name is empty")
	}

	sql, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = sql.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	l.Println("SQLite database is ready")
	return &repo{DB: sql, Logger: l}, nil
}

func (r *repo) Add(user *model.User) (*model.User, error) {
	r.Logger.Println("request add a new user to SQLite database")
	tx := r.DB.Create(&user)
	if tx.Error != nil {
		r.Logger.Printf("Failed adding a new post: %v", tx.Error)
		return nil, tx.Error
	}

	return user, nil
}

func (r *repo) Delete(id int) error {
	r.Logger.Printf("request delete user with ID %v from SQLite database", id)

	var user model.User
	tx := r.DB.Where("id = ?", id).Find(&user)
	if tx.RowsAffected != 0 {
		tx = r.DB.Delete(&user)

		if tx.Error != nil {
			r.Logger.Printf("error while deleting user with ID %v: %v", id, tx.Error)
		} else {
			r.Logger.Printf("user with ID %v has been deleted successfully", id)
		}

		return tx.Error
	}

	err := fmt.Errorf("error: cannot find user with ID %v", id)
	r.Logger.Println(err)

	return err
}

func (r *repo) Get(id int) (*model.User, error) {
	r.Logger.Printf("elaborating the listing request in SQLite database")

	var user *model.User
	tx := r.DB.Where("id = ?", id).Find(&user)

	if tx.RowsAffected != 0 {
		return user, nil
	}

	return nil, fmt.Errorf("user with ID %v not found", id)
}

func (r *repo) GetAll(filters *model.User, pageSize, page int) ([]model.User, error) {
	r.Logger.Printf("elaborating the listing request in SQLite database")

	var users []model.User
	var tx *gorm.DB

	if pageSize > 0 { // pagination has been requested
		if filters != nil { // filtering has been requested
			tx = r.DB.Scopes(paginate(page, pageSize)).Where(&filters).Find(&users)
		} else { // no filtering
			tx = r.DB.Scopes(paginate(page, pageSize)).Find(&users)
		}
	} else { // no pagination
		if filters != nil {
			tx = r.DB.Where(&filters).Find(&users)
		} else {
			tx = r.DB.Find(&users)
		}
	}

	if tx.RowsAffected != 0 {
		r.Logger.Printf("users have been listed successfully from SQLite database")
	} else {
		r.Logger.Printf("there are some problems listing users")
	}

	return users, tx.Error
}

func (r *repo) Update(user, newUser *model.User) (*model.User, error) {
	r.Logger.Printf("elaborating update request in SQLite database")
	tx := r.DB.Model(user).Updates(newUser)

	var err error
	if tx.RowsAffected != 0 {
		r.Logger.Printf("user has been updated successfully in SQLite database")
	} else {
		err = fmt.Errorf("there are some problems updating user with ID %v", user.ID)
		r.Logger.Printf(err.Error())
	}

	return user, err
}
