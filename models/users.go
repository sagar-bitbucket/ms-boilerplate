package models

import "time"

//User Model
type User struct {
	ID        int        `gorm:"id" json:"id"`
	Email     string     `gorm:"email" json:"email"`
	Name      string     `gorm:"name" json:"name"`
	Password  string     `gorm:"password" json:"password"`
	CreatedAt time.Time  `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"-"`
	DeletedAt *time.Time `gorm:"deleted_at" json:"-"`
}

//UserLogin Model
type UserLogin struct {
	Email    string `gorm:"email" json:"email"`
	Password string `gorm:"password" json:"password"`
}
