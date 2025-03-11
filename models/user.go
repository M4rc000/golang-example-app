package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id         int    `gorm:"primaryKey; autoIncrement"`
	Email      string `gorm:"type:varchar(100); unique; not null"`
	Username   string `gorm:"type:varchar(100); unique; null"`
	Name       string `gorm:"type:varchar(100); not null"`
	Picture    string `gorm:"type:varchar(255); null"`
	Created_at string `gorm:"type:varchar(255); not null"`
}
