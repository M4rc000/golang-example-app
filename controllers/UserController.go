package controllers

import (
	"github.com/go-playground/validator/v10"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"golang-example-app/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	URL := strings.Split(c.Request.URL.Path, "/")
	menu := URL[1]
	submenu := URL[2]

	c.HTML(http.StatusFound, "profile.html", gin.H{
		"title":   "User Profile",
		"menu":    helpers.Proper(menu),
		"submenu": helpers.Proper(submenu),
		"user":    userSession,
	})
}

func Index(c *gin.Context) {
	var users []models.User

	config.DB.Where("is_active = ?", 1).Find(&users)

	var DataUsers []map[string]interface{}
	for i, user := range users {
		DataUsers = append(DataUsers, map[string]interface{}{
			"Number":   i + 1,
			"Id":       user.Id,
			"Name":     user.Name,
			"Email":    user.Email,
			"IsActive": user.IsActive,
		})
	}

	// Retrieve flash messages from cookies
	successCreate, _ := c.Cookie("SUCCESS_CREATE")
	failedCreate, _ := c.Cookie("FAILED_CREATE")
	successDelete, _ := c.Cookie("SUCCESS_DELETE")
	failedDelete, _ := c.Cookie("FAILED_DELETE")
	errorMessage, _ := c.Cookie("ERROR")

	// Clear the cookies after retrieving (so the message disappears after refresh)
	c.SetCookie("SUCCESS_CREATE", "", -1, "/", "", false, true)
	c.SetCookie("FAILED_CREATE", "", -1, "/", "", false, true)
	c.SetCookie("SUCCESS_DELETE", "", -1, "/", "", false, true)
	c.SetCookie("FAILED_DELETE", "", -1, "/", "", false, true)
	c.SetCookie("ERROR", "", -1, "/", "", false, true)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":         "User List",
		"users":         DataUsers,
		"successCreate": successCreate,
		"failedCreate":  failedCreate,
		"successDelete": successDelete,
		"failedDelete":  failedDelete,
		"errorMessage":  errorMessage,
	})
}

func CreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{
		"title": "Create User",
	})
}

func StoreUser(c *gin.Context) {
	var validate = validator.New()
	var user models.User

	// Bind request to struct
	if err := c.ShouldBind(&user); err != nil {
		c.SetCookie("ERROR", "Invalid input data", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/users/create")
		return
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	// Validate fields
	if err := validate.Struct(user); err != nil {
		c.SetCookie("ERROR", err.Error(), 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/users/create")
		return
	}

	// Check unique Email on Database
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.SetCookie("ERROR", "Email already exist", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/users/create")
		return
	}

	// Check unique Username on Database
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.SetCookie("ERROR", "Username already exist", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/users/create")
		return
	}

	user.CreatedAt = time.Now()
	user.IsActive = 1

	// Create User
	if err := config.DB.Create(&user).Error; err != nil {
		c.SetCookie("FAILED_CREATE", "Failed to create user", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.SetCookie("SUCCESS_CREATE", "New user successfully created", 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func EditUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	config.DB.First(&user, id)
	c.HTML(http.StatusOK, "edit.html", gin.H{"user": user, "title": "Edit user"})
}

func ShowUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	config.DB.First(&user, id)
	c.HTML(http.StatusOK, "/user/show.html", gin.H{"user": user, "title": "Detail User"})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	config.DB.First(&user, id)

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&user)
	c.Redirect(http.StatusFound, "/")
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	config.DB.Unscoped().Delete(&models.User{}, id)
	c.SetCookie("SUCCESS_DELETE", "User successfully deleted", 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}
