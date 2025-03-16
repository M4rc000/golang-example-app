package config

import (
	"fmt"
	"golang-example-app/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Migration failed")
	} else {
		fmt.Println("Migration Successfully")
	}
}
