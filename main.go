package main

import (
	"golang-example-app/config"
)

func main() {
	config.ConnectionDB()
	defer config.CloseDB()
	//routes.RunApp()
	config.Migrate()
}
