package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-example-app/config"
	"golang-example-app/models"
	"net/http"
)

func Register(c *gin.Context) {
	session := sessions.Default(c)

	// RETRIEVE FLASH MESSAGES
	failedRegister := session.Get("FAILED_REGISTER")
	errorInputData := session.Get("ERROR_INPUTDATA")
	errorUsername := session.Get("ERROR_USERNAME")
	errorName := session.Get("ERROR_NAME")
	errorEmail := session.Get("ERROR_EMAIL")
	errorPassword := session.Get("ERROR_PASSWORD")
	duplicateEmail := session.Get("DUPLICATE_EMAIL")
	duplicateUsername := session.Get("DUPLICATE_USERNAME")

	//Clear flash messages (so the message disappears after refresh)
	session.Delete("FAILED_REGISTER")
	session.Delete("ERROR")
	session.Delete("ERROR_USERNAME")
	session.Delete("ERROR_NAME")
	session.Delete("ERROR_EMAIL")
	session.Delete("ERROR_PASSWORD")
	session.Delete("DUPLICATE_EMAIL")
	session.Delete("DUPLICATE_USERNAME")

	session.Save()

	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":             "Register",
		"failedRegister":    failedRegister,
		"errorUsername":     errorUsername,
		"errorName":         errorName,
		"errorEmail":        errorEmail,
		"errorPassword":     errorPassword,
		"errorInputData":    errorInputData,
		"duplicateEmail":    duplicateEmail,
		"duplicateUsername": duplicateUsername,
		"old": gin.H{
			"Name":     c.DefaultQuery("Name", ""),
			"Username": c.DefaultQuery("Username", ""),
			"Email":    c.DefaultQuery("Email", ""),
		},
	})
}

func StoreRegister(c *gin.Context) {
	session := sessions.Default(c)
	var validate = validator.New()
	var user models.User
	var errors = make(map[string]string)

	// Use ShouldBindJSON if receiving JSON requests
	if err := c.ShouldBind(&user); err != nil {
		session.Set("ERROR_INPUTDATA", "Invalid input data")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
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
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Hash password before saving
	if err := user.HashPassword(); err != nil {
		session.Set("ERROR", "Failed to hash password")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
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
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Activate the user
	user.IsActive = 1

	// Save user to database
	if err := config.DB.Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", "Failed to create user")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	session.Set("SUCCESS_REGISTER", "Registration successful")
	session.Save()
	c.Redirect(http.StatusFound, "/auth/")
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	successRegister := session.Get("SUCCESS_REGISTER")
	loginError := session.Get("LOGIN_ERROR")

	session.Delete("LOGIN_ERROR")
	session.Delete("SUCCESS_REGISTER")

	session.Save()

	c.HTML(http.StatusFound, "login.html", gin.H{
		"title":           "Login",
		"successRegister": successRegister,
		"loginError":      loginError,
	})
}

func Authenticate(c *gin.Context) {
	session := sessions.Default(c)

	var request struct {
		EmailUsername string `form:"email-username" binding:"required"`
		Password      string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		session.Set("LOGIN_ERROR", "Invalid email/username or password.")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/")
		return
	}

	var user models.User
	result := config.DB.Where("email = ? OR username = ?", request.EmailUsername, request.EmailUsername).First(&user)
	if result.Error != nil {
		session.Set("LOGIN_ERROR", "Invalid email/username or password.")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/")
		return
	}

	if !user.CheckPassword(request.Password) {
		session.Set("LOGIN_ERROR", "Invalid email/username or password.")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/")
		return
	}

	// Authentication successful
	session.Set("USER_ID", user.ID) // Store user ID in session
	session.Save()
	c.Redirect(http.StatusFound, "/home/dashboard")
}
