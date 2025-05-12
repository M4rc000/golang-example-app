package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id         int    `gorm:"primaryKey; autoIncrement"`
	Email      string `gorm:"type:varchar(100); unique; not null" form:"Email" json:"Email" validate:"required,email,max=100"`
	Username   string `gorm:"type:varchar(100); unique; not null" form:"Username" json:"Username" validate:"required,max=100"`
	Name       string `gorm:"type:varchar(100); not null" form:"Name" json:"Name" validate:"required,max=100"`
	Password   string `gorm:"type:varchar(255); not null" form:"Password" json:"Password" validate:"required,min=4"`
	Gender     string `gorm:"type:varchar(10); null" form:"Gender" json:"Gender"`
	Address    string `gorm:"type:varchar(255); null" form:"Address" json:"Address"`
	PostalCode string `gorm:"type:varchar(10); null" form:"PostalCode" json:"PostalCode"`
	Country    string `gorm:"type:varchar(100); null" form:"Country" json:"Country"`
	Picture    string `gorm:"type:varchar(255); default:profile-user/default.png" form:"Picture" json:"Picture"`
	Token      string `gorm:"type:varchar(100)" form:"Token" json:"Token"`
	RoleID     int    `gorm:"type: int" form:"RoleID" json:"RoleID"`
	IsActive   int    `gorm:"type:int;default:1"`
}

type UserProfile struct {
	gorm.Model
	Id         int    `gorm:"primaryKey; autoIncrement"`
	Email      string `gorm:"type:varchar(100); unique; not null" form:"Email" json:"Email" validate:"required,email,max=100"`
	Username   string `gorm:"type:varchar(100); unique; not null" form:"Username" json:"Username" validate:"required,max=100"`
	Name       string `gorm:"type:varchar(100); not null" form:"Name" json:"Name" validate:"required,max=100"`
	Gender     string `gorm:"type:varchar(10); null" form:"Gender" json:"Gender"`
	Address    string `gorm:"type:varchar(255); null" form:"Address" json:"Address"`
	PostalCode string `gorm:"type:varchar(10); null" form:"PostalCode" json:"PostalCode"`
	Country    string `gorm:"type:varchar(100); null" form:"Country" json:"Country"`
	Picture    string `gorm:"type:varchar(255); default:profile-user/default.png" form:"-" json:"-"`
	IsActive   int    `gorm:"type:int;default:0"`
}

type AddManageUser struct {
	gorm.Model
	Id       int    `gorm:"primaryKey; autoIncrement"`
	Email    string `gorm:"type:varchar(100); unique; not null" form:"Email" validate:"required,email,max=100"`
	Username string `gorm:"type:varchar(100); unique; not null" form:"Username" validate:"required,max=100"`
	Name     string `gorm:"type:varchar(100); not null" form:"Name" validate:"required,max=100"`
	Password string `gorm:"type:varchar(255); not null" form:"Password" validate:"required,min=4"`
	Gender   string `gorm:"type:varchar(10); null" form:"Gender" validate:"required"`
	IsActive int    `gorm:"type:int;default:1"`
}

type UpdateManageUser struct {
	gorm.Model
	Id       int    `gorm:"primaryKey; autoIncrement"`
	Email    string `gorm:"type:varchar(100); unique; not null" form:"Email" validate:"required,email,max=100"`
	Username string `gorm:"type:varchar(100); unique; not null" form:"Username" validate:"required,max=100"`
	Name     string `gorm:"type:varchar(100); not null" form:"Name" validate:"required,max=100"`
	Password string `form:"Password" validate:"omitempty,min=4"`
	Gender   string `gorm:"type:varchar(10); null" form:"Gender" validate:"required"`
	IsActive int    `gorm:"type:int;default:1"`
}
