package models

import (
	"final-project/helpers"

	"gorm.io/gorm"
)

// User represents the model for an user
type User struct {
	GORMModel
	Username string `gorm:"not null;uniqueIndex" json:"username" valid:"required~Username is required"`
	Email    string `gorm:"not null;uniqueIndex" json:"email" valid:"required~Email is required, email~Invalid email format"`
	Password string `gorm:"not null" json:"password" valid:"required~Password is required, minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Age      int    `gorm:"not null" json:"age" valid:"required~Age is required ,range(8|500)~You must be at least 8 years old"`
	// Comments    []Comment     `json:"comments"`
	// Photos      []Photo       `json:"photos"`
	// SocialMedia []SocialMedia `json:"social_media"`
}

// Run Autommatically before Creating (hashing password)
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	// Call helpers.HashPassword to HASH Password
	hashedPass, err := helpers.HashPassword(u.Password)
	if err != nil {
		return
	}

	// change Password to Hashed Password
	u.Password = hashedPass

	return
}
