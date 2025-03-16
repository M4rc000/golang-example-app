package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       int    `gorm:"primaryKey; autoIncrement"`
	Email    string `gorm:"type:varchar(100); unique; not null" validate:"required,email,max=100"`
	Username string `gorm:"type:varchar(100); unique; not null" validate:"required,max=100"`
	Name     string `gorm:"type:varchar(100); not null" validate:"required,max=100"`
	Password string `gorm:"type:varchar(255); not null" validate:"required,min=4"` // Store hashed password
	Gender   string `gorm:"type:varchar(10); null"`
	Picture  string `gorm:"type:varchar(255); null"`
	IsActive int    `gorm:"type:int;default:0"`
}

// HashPassword hashes the user's password before saving to the database
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password is correct
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
