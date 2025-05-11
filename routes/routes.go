package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"golang-example-app/controllers"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"html/template"
	"log"
	"net/http"
	"os"
)

func RunApp() {
	app := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   3600,                 // Session lifetime in seconds (1 hour)
		Path:     "/",                  // Cookie available on all paths
		HttpOnly: true,                 // Prevent access from JavaScript (recommended)
		Secure:   false,                // Set to true in production if using HTTPS
		SameSite: http.SameSiteLaxMode, // Controls cross-site cookie behavior
	})
	app.Use(sessions.Sessions("prime-app", store))

	// CRSF TOKEN
	app.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CRSF_TOKEN"),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	app.Static("/assets", "./assets")

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	}).ParseGlob("views/**/*"))

	app.SetHTMLTemplate(tmpl)

	// HANDLE NOT FOUND
	app.NoRoute(controllers.NoFoundRoute)

	app.GET("/", helpers.RedirectSlashRoute("auth"))

	// AUTHENTICATION
	authGroup := app.Group("/auth")
	{
		authGroup.Use(middlewares.GuestRequired)
		authGroup.GET("", controllers.Login)
		authGroup.POST("/login", controllers.Authenticate)

		// Move register routes inside auth group
		authGroup.GET("/register", controllers.Register)
		authGroup.POST("/register", controllers.StoreRegister)
	}

	// ADMINISTRATOR
	administratorGroup := app.Group("/administrator")
	{
		administratorGroup.Use(middlewares.AuthRequired)
		administratorGroup.GET("", helpers.RedirectSlashRoute("/administrator/manage-user"))
		administratorGroup.GET("/manage-user", controllers.ManageUser)
		administratorGroup.GET("/add-new-user", controllers.AddNewUser)
		administratorGroup.POST("/save-user", controllers.SaveNewUser)
		administratorGroup.GET("/edit-user", controllers.EditUser)
		administratorGroup.POST("/update-user", controllers.UpdateUser)
		administratorGroup.GET("/show-user/:id", controllers.ShowUser)
		administratorGroup.GET("/manage-role", controllers.ManageRole)
		administratorGroup.GET("/manage-menu", controllers.ManageMenu)
		administratorGroup.GET("/manage-submenu", controllers.ManageSubmenu)
	}

	// DASHBOARD
	dashboardGroup := app.Group("/dashboard")
	{
		dashboardGroup.Use(middlewares.AuthRequired)
		dashboardGroup.GET("", controllers.Dashboard)
	}

	// HR
	hrGroup := app.Group("/hr")
	{
		hrGroup.Use(middlewares.AuthRequired)
		hrGroup.GET("", helpers.RedirectSlashRoute("/employees"))
		hrGroup.GET("/employees", controllers.Employees)
		hrGroup.GET("/attendance", controllers.Attendance)
	}

	// USER
	userGroup := app.Group("/user")
	{
		userGroup.Use(middlewares.AuthRequired)
		userGroup.GET("/profile", controllers.UserProfile)
		userGroup.POST("/profile", controllers.UpdateUserProfile)
		userGroup.GET("/settings", controllers.CreateUser)
		userGroup.GET("/logout", controllers.Logout)
	}

	err := app.Run(":5000")
	if err != nil {
		log.Fatal("Application failed to run")
	}
	log.Println("App running on port 5000")
}
