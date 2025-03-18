package models

type Movie struct {
	ID          int    `gorm:"primaryKey; autoIncrement"`
	Title       string `gorm:"type:varchar(50); unique; not null" validate:"required,max=50"`
	Description string `gorm:"type:text; not null" validate:"required"`
	Year        string `gorm:"type:varchar(10); not null" validate:"required"`
	Duration    string `gorm:"type:varchar(30); not null" validate:"required"`
	Age         string `gorm:"type:varchar(20); not null" validate:"required"`
	Rating      int    `gorm:"type:int; default:0"`
	Picture     string `gorm:"type:varchar(255); not null" validate:"required"`
	IsActive    int    `gorm:"type:int;default:1"`
}
