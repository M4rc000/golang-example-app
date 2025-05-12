package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Id       int    `gorm:"primaryKey; autoIncrement"`
	Name     string `gorm:"type:varchar(50); unique; not null" form:"Name" json:"Name" validate:"required,max=100"`
	IsActive int    `gorm:"type:int;default:1"`
}
