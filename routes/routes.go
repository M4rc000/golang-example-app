package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"golang-example-app/controllers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunApp() {
	app := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	app.Use(sessions.Sessions("mysession", store))
	store.Options(sessions.Options{
		MaxAge: 3600, // 1 hour
	})

	app.Static("/assets", "./assets")
	app.LoadHTMLGlob("views/**/*")

	// AUTHENTICATION
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth/")
	})
	app.GET("/auth/", controllers.Login)
	app.POST("/auth/", controllers.Authenticate)
	app.GET("/auth/register", controllers.Register)
	app.POST("/auth/register", controllers.StoreRegister)
	//app.GET("/auth/forgot-password", controllers.ForgotPassword)

	//app.GET("/", controllers.Index)
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
