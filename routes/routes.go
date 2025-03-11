package routes

import (
	"golang-example-app/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func RunApp() {
	app := gin.Default()
	app.LoadHTMLGlob("views/**/*")

	app.GET("/users", controllers.GetUsers)
	// app.POST("/", controllers.CreateUser)
	// app.GET("/edit/:id", controllers.EditUserForm)
	// app.POST("/update/:id", controllers.UpdateUser)
	// app.GET("/delete/:id", controllers.DeleteUser)

	err := app.Run(":8080")
	if err != nil {
		log.Fatal("Application failed to run")
	}
}
