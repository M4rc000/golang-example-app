package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"golang-example-app/models"
	"net/http"
)

func ManageUser(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	helpers.FlashMessage(c, "USER_EMAIL_USERNAME")
	var users []models.User

	config.DB.Where("is_active = ?", 1).Find(&users)

	var DataUsers []map[string]interface{}
	for i, user := range users {
		DataUsers = append(DataUsers, map[string]interface{}{
			"Number":   i + 1,
			"Id":       user.Id,
			"Picture":  user.Picture,
			"Name":     user.Name,
			"Username": user.Username,
			"Email":    user.Email,
			"Gender":   user.Gender,
			"IsActive": user.IsActive,
		})
	}
	c.HTML(http.StatusOK, "manage_user.html", gin.H{
		"title":     "Manage User",
		"menus":     menus,
		"user":      userSession,
		"DataUsers": DataUsers,
	})
}

func AddNewUser(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, submenu := helpers.GetMenuSubmenu(c)

	c.HTML(http.StatusFound, "add_new_user.html", gin.H{
		"title":     "Create New User",
		"csrfToken": csrf.GetToken(c),
		"menu":      menu,
		"submenu":   submenu,
		"user":      userSession,
	})
}

func SaveNewUser(c *gin.Context) {
	session := sessions.Default(c)
	var validate = validator.New()
	var user models.User
	var errors = make(map[string]string)

	// Use ShouldBindJSON if receiving JSON requests
	if err := c.ShouldBind(&user); err != nil {
		session.Set("ERROR_INPUTDATA", "Invalid input data")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage_user")
		return
	}

	// Validate user struct
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errors["ERROR_EMAIL"] = "Email must be a valid email address"
			case "Username":
				errors["ERROR_USERNAME"] = "Username must be a valid username"
			case "Name":
				errors["ERROR_NAME"] = "Name is required"
			case "Password":
				errors["ERROR_PASSWORD"] = "Password must be at least 8 characters!"
			}
		}

		// Store errors in cookies
		for key, msg := range errors {
			session.Set(key, msg)
			err := session.Save()
			if err != nil {
				c.JSON(http.StatusFound, gin.H{
					"error": err.Error(),
				})
			}
		}

		// Redirect once after storing errors
		c.Redirect(http.StatusFound, "/administrator/manage_user")
		return
	}

	// Hash password before saving
	if err := user.HashPassword(); err != nil {
		session.Set("ERROR", "Failed to hash password")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage_user")
		return
	}

	// Check for duplicate email or username in a single query
	var existingUser models.User
	if err := config.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err == nil {
		if existingUser.Email == user.Email {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
		}
		if existingUser.Username == user.Username {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage_user")
		return
	}

	// Activate the user
	user.IsActive = 1

	// Save user to database
	if err := config.DB.Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", "Failed to create user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage_user")
		return
	}

	session.Set("SUCCESS_REGISTER", "Registration successful")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage_user")
}

func ManageRole(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "manage_role.html", gin.H{
		"title": "Manage Role",
		"menus": menus,
		"user":  userSession,
	})
}

func ManageMenu(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "manage_menu.html", gin.H{
		"title": "Manage Menu",
		"menus": menus,
		"user":  userSession,
	})
}

func ManageSubmenu(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "manage_submenu.html", gin.H{
		"title": "Manage Submenu",
		"menus": menus,
		"user":  userSession,
	})
}
