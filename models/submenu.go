package models

import "gorm.io/gorm"

type SubMenu struct {
	gorm.Model
	MenuID   uint   `gorm:"not null" validate:"required"`
	Name     string `gorm:"type:varchar(50);not null" validate:"required,max=50"`
	URL      string `gorm:"type:varchar(50);not null" validate:"required,max=50"`
	Icon     string `gorm:"type:varchar(50);not null" validate:"required,max=50"`
	IsActive int    `gorm:"default:1"`
}
