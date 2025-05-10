package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Id       int       `gorm:"primaryKey; autoIncrement"`
	Name     string    `gorm:"type:varchar(100); unique; not null" validate:"required,max=20"`
	SubMenus []SubMenu `gorm:"foreignKey:MenuID"`
}
