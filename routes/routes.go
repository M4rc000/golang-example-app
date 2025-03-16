package routes

import (
	"golang-example-app/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func RunApp() {
	app := gin.Default()
	app.LoadHTMLGlob("views/*.html")

	app.GET("/", controllers.Index)
	app.GET("/users/create", controllers.CreateUser)
	app.POST("/users/create", controllers.StoreUser)
	app.GET("/users/edit/:id", controllers.EditUser)
	app.GET("/users/show/:id", controllers.ShowUser)
	app.POST("/users/update/:id", controllers.UpdateUser)
	app.GET("/users/delete/:id", controllers.DeleteUser)

	err := app.Run(":8080")
	if err != nil {
		log.Fatal("Application failed to run")
	}
}
