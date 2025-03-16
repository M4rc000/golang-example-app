package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-example-app/config"
	"golang-example-app/models"
	"net/http"
)

func Register(c *gin.Context) {
	// Retrieve flash messages from cookies
	successCreate, _ := c.Cookie("SUCCESS_CREATE")
	failedCreate, _ := c.Cookie("FAILED_CREATE")
	errorMessage, _ := c.Cookie("ERROR")

	// Clear the cookies after retrieving (so the message disappears after refresh)
	c.SetCookie("FAILED_CREATE", "", -1, "/auth/register", "", false, true)
	c.SetCookie("ERROR", "", -1, "/auth/register", "", false, true)

	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":         "Register",
		"successCreate": successCreate,
		"failedCreate":  failedCreate,
		"errorMessage":  errorMessage,
	})
}

func StoreRegister(c *gin.Context) {
	var validate = validator.New()
	var user models.User

	// Use ShouldBindJSON if receiving JSON requests
	if err := c.ShouldBind(&user); err != nil {
		c.SetCookie("ERROR", "Invalid input data", 3600, "/auth/register", "", false, false)
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Validate user struct
	if err := validate.Struct(user); err != nil {
		c.SetCookie("ERROR", err.Error(), 3600, "/auth/register", "", false, false)
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Hash password before saving
	if err := user.HashPassword(); err != nil {
		c.SetCookie("ERROR", "Failed to hash password", 3600, "/auth/register", "", false, false)
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Check unique email
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.SetCookie("ERROR", "Email already exists", 3600, "/auth/register", "", false, false)
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Check unique username
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.SetCookie("ERROR", "Username already exists", 3600, "/auth/register", "", false, false)
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Activate the user
	user.IsActive = 1

	// Save user to database
	if err := config.DB.Create(&user).Error; err != nil {
		c.SetCookie("FAILED_REGISTER", "Failed to create user", 3600, "/auth/register", "", false, false)
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	c.SetCookie("SUCCESS_REGISTER", "Registration successful", 3600, "/auth", "", false, false)
	c.Redirect(http.StatusFound, "/auth")
}

func Login(c *gin.Context) {
	successCreate, _ := c.Cookie("SUCCESS_CREATE")
	c.SetCookie("SUCCESS_CREATE", "", -1, "/auth/", "", false, true)
	c.HTML(http.StatusFound, "login.html", gin.H{
		"title":         "Login",
		"successCreate": successCreate,
	})
}

func Authenticate(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if !user.CheckPassword(request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// TODO: Generate JWT token (future step)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
