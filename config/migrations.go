package config

import (
	"fmt"
	"golang-example-app/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Migration User failed", err)
		return
	}
	err = DB.AutoMigrate(&models.Movie{})
	if err != nil {
		fmt.Println("Migration Movie failed", err)
		return
	}
	err = DB.AutoMigrate(&models.Menu{})
	if err != nil {
		fmt.Println("Migration Menu failed", err)
		return
	}
	err = DB.AutoMigrate(&models.SubMenu{})
	if err != nil {
		fmt.Println("Migration Submenu failed", err)
		return
	}
}
