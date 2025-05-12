package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/models"
	"net/http"
)

func Register(c *gin.Context) {
	err := helpers.FlashMessage(c, "ERROR")
	failedRegister := helpers.FlashMessage(c, "FAILED_REGISTER")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	errorUsername := helpers.FlashMessage(c, "ERROR_USERNAME")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorEmail := helpers.FlashMessage(c, "ERROR_EMAIL")
	errorPassword := helpers.FlashMessage(c, "ERROR_PASSWORD")
	duplicateEmail := helpers.FlashMessage(c, "DUPLICATE_EMAIL")
	duplicateUsername := helpers.FlashMessage(c, "DUPLICATE_USERNAME")

	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":             "Register",
		"csrfToken":         csrf.GetToken(c),
		"failedRegister":    failedRegister,
		"err":               err,
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
				session.Set("ERROR_INPUTDATA", err.Error())
				session.Save()
			}
		}

		// Redirect once after storing errors
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Hash password before saving
	password, err := helpers.HashPassword(user.Password)
	if err != nil {
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
	user.Password = password

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
	successRegister := helpers.FlashMessage(c, "SUCCESS_REGISTER")
	loginError := helpers.FlashMessage(c, "LOGIN_ERROR")

	c.HTML(http.StatusFound, "login.html", gin.H{
		"title":           "Login",
		"csrfToken":       csrf.GetToken(c),
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
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	fmt.Printf("Username: %s, Password: %s", request.EmailUsername, request.Password)

	var user models.User
	result := config.DB.Where("email = ? OR username = ?", request.EmailUsername, request.EmailUsername).First(&user)
	if result.Error != nil {
		session.Set("LOGIN_ERROR", "Invalid email/username or password.")
		session.Save()
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	password := helpers.CheckPasswordHash(request.Password, user.Password)
	if !password {
		session.Set("LOGIN_ERROR", "Invalid email/username or password.")
		session.Save()
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	userID := user.Id
	userEmail := user.Email
	userUsername := user.Username
	session.Set("USER_ID", userID)
	session.Set("USER_EMAIL", userEmail)
	session.Set("USER_USERNAME", userUsername)
	session.Save()

	// Authentication successful, Locate user based on access menu
	firstAccessMenu := helpers.GetFirstAccessibleURL(config.DB, user.RoleID)
	c.Redirect(http.StatusFound, firstAccessMenu)
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/auth/")
}
