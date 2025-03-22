package config

import (
	"fmt"
	"golang-example-app/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Migration failed", err)
		return
	}
	err = DB.AutoMigrate(&models.Movie{})
	if err != nil {
		fmt.Println("Migration failed", err)
		return
	}
}
