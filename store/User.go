package store

import "time"

type User struct {
	Id        uint `gorm:"primaryKey"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Nickname  string `json:"nickname" gorm:"not null"`
	Password  string `json:"password" gorm:"not null"`
	Email     string `json:"email" gorm:"not null"`
	Country   string `json:"country" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
