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

	// HANDLE NOT FOUND
	app.NoRoute(controllers.NoFoundRoute)

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth")
	})

	// AUTHENTICATION
	// AUTHENTICATION
	authGroup := app.Group("/auth")
	{
		authGroup.Use(middlewares.GuestRequired)
		authGroup.GET("", controllers.Login)
		authGroup.POST("", controllers.Authenticate)

		// Move register routes inside auth group
		authGroup.GET("/register", controllers.Register)
		authGroup.POST("/register", controllers.StoreRegister)
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
		userGroup.POST("/profile", controllers.UpdateUserProfile)
		userGroup.GET("/settings", controllers.CreateUser)
	}

	err := app.Run(":3000")
	if err != nil {
		log.Fatal("Application failed to run")
	}
}
