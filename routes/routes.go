package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang-example-app/controllers"
	"golang-example-app/middlewares"
	"log"
	"net/http"
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

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth")
	})

	// AUTHENTICATION
	authGroup := app.Group("/auth")
	{
		authGroup.Use(middlewares.GuestRequired)
		authGroup.GET("", controllers.Login)
		authGroup.POST("", controllers.Authenticate)
		app.GET("/register", controllers.Register)
		app.POST("/register", controllers.StoreRegister)
	}

	// HOME
	homeGroup := app.Group("/home")
	{
		homeGroup.Use(middlewares.AuthRequired)
		homeGroup.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/home/dashboard")
		})
		homeGroup.GET("/dashboard", controllers.Dashboard)
		homeGroup.GET("/logout", controllers.Logout)
	}

	userGroup := app.Group("/user")
	{
		userGroup.Use(middlewares.AuthRequired)
		userGroup.GET("/profile", controllers.UserProfile)
		userGroup.GET("/settings", controllers.CreateUser)

		//userGroup.GET("/users/create", controllers.CreateUser)
		//userGroup.POST("/users/create", controllers.StoreUser)
		//userGroup.GET("/users/edit/:id", controllers.EditUser)
		//userGroup.GET("/users/show/:id", controllers.ShowUser)
		//userGroup.POST("/users/update/:id", controllers.UpdateUser)
		//userGroup.GET("/users/delete/:id", controllers.DeleteUser)
	}

	err := app.Run(":8000")
	if err != nil {
		log.Fatal("Application failed to run")
	}
}
