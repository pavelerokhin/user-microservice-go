package model

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id" bson:"id"`
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	Nickname  string    `json:"nickname" bson:"nickname"`
	Password  string    `json:"password" bson:"password"`
	Email     string    `json:"email" bson:"email"`
	Country   string    `json:"country" bson:"country"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
