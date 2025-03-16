package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       int    `gorm:"primaryKey; autoIncrement"`
	Email    string `gorm:"type:varchar(100); unique; not null" validate:"required,email,max=100"`
	Username string `gorm:"type:varchar(100); unique; not null" validate:"required,max=100"`
	Name     string `gorm:"type:varchar(100); not null" validate:"required,max=100"`
	Gender   string `gorm:"type:varchar(10); null"`
	Picture  string `gorm:"type:varchar(255); null"`
	IsActive int    `gorm:"type:int;default:0"`
}
