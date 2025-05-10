package models

import "gorm.io/gorm"

type UserAccessMenu struct {
	gorm.Model
	Id     int `gorm:"primaryKey; autoIncrement;"`
	MenuID int `gorm:"not null" validate:"required,max=2"`
	RoleID int `gorm:"not null; validate: required,max=2"`
}
