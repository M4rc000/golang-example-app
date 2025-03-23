package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"golang-example-app/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UserProfile(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, submenu := helpers.GetMenuSubmenu(c)

	session := sessions.Default(c)
	errorUpdateProfile := session.Get("ERROR_UPDATEPROFILE")
	successUpdateProfile := session.Get("SUCCESS_UPDATEPROFILE")
	failedUpdateProfile := session.Get("FAILED_UPDATEPROFILE")
	errorUsername := session.Get("ERROR_USERNAME")
	errorName := session.Get("ERROR_NAME")
	errorEmail := session.Get("ERROR_EMAIL")
	errorPicture := session.Get("ERROR_PICTURE")
	duplicateEmail := session.Get("DUPLICATE_EMAIL")
	duplicateUsername := session.Get("DUPLICATE_USERNAME")

	session.Delete("ERROR_PICTURE")
	session.Delete("ERROR_UPDATEPROFILE")
	session.Delete("SUCCESS_UPDATEPROFILE")
	session.Delete("FAILED_UPDATEPROFILE")
	session.Delete("ERROR_USERNAME")
	session.Delete("ERROR_NAME")
	session.Delete("ERROR_EMAIL")
	session.Delete("DUPLICATE_EMAIL")
	session.Delete("DUPLICATE_USERNAME")

	session.Save()

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title":                "User Profile",
		"menu":                 menu,
		"submenu":              submenu,
		"user":                 userSession,
		"successUpdateProfile": successUpdateProfile,
		"failedUpdateProfile":  failedUpdateProfile,
		"errorUsername":        errorUsername,
		"errorName":            errorName,
		"errorPicture":         errorPicture,
		"error":                errorUpdateProfile,
		"errorEmail":           errorEmail,
		"duplicateEmail":       duplicateEmail,
		"duplicateUsername":    duplicateUsername,
	})
}

func UpdateUserProfile(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	session := sessions.Default(c)
	var validate = validator.New()
	var user models.UserProfile
	var errors = make(map[string]string)

	if err := c.ShouldBind(&user); err != nil {
		log.Fatal(err)
		session.Set("ERROR_UPDATEPROFILE", "Invalid input data")
		session.Save()
		c.Redirect(http.StatusFound, "/user/profile")
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
		c.Redirect(http.StatusFound, "/user/profile")
		return
	}

	// FILE UPLOAD
	file, err := c.FormFile("Picture")
	picturePath := ""
	if err == nil && file != nil {
		// Validate file type.
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
			"image/webp": true,
		}

		if !allowedTypes[file.Header.Get("Content-Type")] {
			session.Set("ERROR_PICTURE", "Invalid file format. Only Webp, JPG, and PNG are allowed.")
			session.Save()
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}

		if file.Size > 1048576 {
			session.Set("ERROR_PICTURE", "File size exceeds maximum allowed limit of 1 MB")
			session.Save()
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}

		// Generate a new unique filename.
		fileExt := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%d%s", userSession.Id, fileExt)

		// Define the destination directory.
		saveDir := "assets/img/profile-user"

		// Create the directory if it doesn't exist.
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			log.Println("Error creating directory:", err)
			session.Set("ERROR_PICTURE", "Failed to create directory for profile picture")
			session.Save()
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}

		// Create the full file destination path.
		destinationPath := filepath.Join(saveDir, newFilename)
		if err := c.SaveUploadedFile(file, destinationPath); err != nil {
			log.Println("Error saving file:", err)
			session.Set("ERROR_PICTURE", "Failed to upload profile picture")
			session.Save()
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}

		// Set the picturePath for the DB update (relative to your assets folder).
		picturePath = filepath.ToSlash(filepath.Join("profile-user", newFilename))
	}

	// CHECK DUPLICATE EMAIL OR USERNAME
	var existingUser models.UserProfile

	// Check username duplication (excluding current user)
	if userSession.Username != user.Username {
		if err := config.DB.Where("username = ? AND id != ?", user.Username, userSession.Id).First(&existingUser).Error; err == nil {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
			session.Save()
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}
	}

	// Check email duplication (excluding current user)
	if userSession.Email != user.Email {
		if err := config.DB.Where("email = ? AND id != ?", user.Email, userSession.Id).First(&existingUser).Error; err == nil {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
			session.Save()
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}
	}

	// Prepare the update struct. For fields that are not changing, you can retain the old values.
	updates := models.UserProfile{
		Name:       user.Name,
		Email:      user.Email,
		Username:   user.Username,
		Gender:     user.Gender,
		Address:    user.Address,
		PostalCode: user.PostalCode,
		Country:    user.Country,
	}
	if picturePath != "" {
		updates.Picture = picturePath
	}

	// Save Data Profile to Database
	userID := userSession.Id
	result := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&updates)
	if result.Error != nil {
		log.Println("Update error:", result.Error)
		session.Set("FAILED_UPDATEPROFILE", "Failed to update user profile")
		session.Save()
		c.Redirect(http.StatusFound, "/user/profile")
		return
	}
	if result.RowsAffected == 0 {
		session.Set("FAILED_UPDATEPROFILE", "No changes were made to the profile.")
		session.Save()
		c.Redirect(http.StatusFound, "/user/profile")
		return
	}

	session.Set("SUCCESS_UPDATEPROFILE", "Profile successfully updated")
	session.Save()
	c.Redirect(http.StatusFound, "/user/profile")
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
