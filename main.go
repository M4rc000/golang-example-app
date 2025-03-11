package main

import (
	"golang-example-app/config"
	"golang-example-app/routes"
)

func main() {
	config.ConnectionDB()
	routes.RunApp()
}
