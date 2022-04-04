package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	DB *gorm.DB
}

func NewSQLite() (*DB, error) {
	var db DB
	sqlite, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.DB = sqlite

	err = sqlite.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}
	return &db, nil
}
