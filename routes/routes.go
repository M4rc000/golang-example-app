package routes

import (
	"golang-example-app/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func RunApp() {
	app := gin.Default()
	app.Static("/assets", "./assets")
	app.LoadHTMLGlob("views/**/*")

	// AUTHENTICATION
	app.GET("/auth/", controllers.Login)
	app.GET("/auth/register", controllers.Register)
	app.POST("/auth/register", controllers.StoreRegister)
	//app.GET("/auth/forgot-password", controllers.ForgotPassword)

	app.GET("/", controllers.Index)
	app.GET("/users/create", controllers.CreateUser)
	app.POST("/users/create", controllers.StoreUser)
	app.GET("/users/edit/:id", controllers.EditUser)
	app.GET("/users/show/:id", controllers.ShowUser)
	app.POST("/users/update/:id", controllers.UpdateUser)
	app.GET("/users/delete/:id", controllers.DeleteUser)

	err := app.Run(":8000")
	if err != nil {
		log.Fatal("Application failed to run")
	}
}
