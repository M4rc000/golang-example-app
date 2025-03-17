package main

import (
	"golang-example-app/config"
	"golang-example-app/routes"
)

func main() {
	config.ConnectionDB()
	defer config.CloseDB()
	routes.RunApp()
	//config.Migrate()
}
