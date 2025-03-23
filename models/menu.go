package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name string `gorm:"type:varchar(100); unique; not null" validate:"required,max=20"`
}
