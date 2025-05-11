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
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	err := helpers.FlashMessage(c, "ERROR")
	successRegister := helpers.FlashMessage(c, "SUCCESS_REGISTER")

	var users []models.User
	config.DB.Where("is_active = ?", 1).Find(&users)

	var DataUsers []map[string]interface{}
	for i, user := range users {
		DataUsers = append(DataUsers, map[string]interface{}{
			"Number":   i + 1,
			"Id":       helpers.EncodeID(user.Id),
			"Picture":  user.Picture,
			"Name":     user.Name,
			"Username": user.Username,
			"Email":    user.Email,
			"Gender":   user.Gender,
			"IsActive": user.IsActive,
		})
	}

	c.HTML(http.StatusOK, "manage_user.html", gin.H{
		"title":           "Manage User",
		"menus":           menus,
		"menu":            menu,
		"user":            userSession,
		"DataUsers":       DataUsers,
		"csrfToken":       csrf.GetToken(c),
		"err":             err,
		"successRegister": successRegister,
	})
}

func AddNewUser(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	err := helpers.FlashMessage(c, "ERROR")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	duplicateEmail := helpers.FlashMessage(c, "DUPLICATE_EMAIL")
	duplicateUsername := helpers.FlashMessage(c, "DUPLICATE_USERNAME")
	failedRegister := helpers.FlashMessage(c, "FAILED_REGISTER")
	errorEmail := helpers.FlashMessage(c, "ERROR_EMAIL")
	errorUsername := helpers.FlashMessage(c, "ERROR_USERNAME")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorPassword := helpers.FlashMessage(c, "ERROR_PASSWORD")
	errorGender := helpers.FlashMessage(c, "ERROR_GENDER")

	c.HTML(http.StatusOK, "add_new_user.html", gin.H{
		"title":             "New User",
		"menu":              menu,
		"menus":             menus,
		"user":              userSession,
		"csrfToken":         csrf.GetToken(c),
		"err":               err,
		"errorInputData":    errorInputData,
		"duplicateUsername": duplicateUsername,
		"duplicateEmail":    duplicateEmail,
		"failedRegister":    failedRegister,
		"errorEmail":        errorEmail,
		"errorUsername":     errorUsername,
		"errorName":         errorName,
		"errorPassword":     errorPassword,
		"errorGender":       errorGender,
	})
}

func SaveNewUser(c *gin.Context) {
	session := sessions.Default(c)
	var validate = validator.New()
	var user models.AddManageUser
	var errors = make(map[string]string)

	if err := c.ShouldBind(&user); err != nil {
		session.Set("ERROR_INPUTDATA", "Invalid input data")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				if err.Tag() == "required" {
					errors["ERROR_EMAIL"] = "Email must be a valid email address"
				} else if err.Tag() == "email" {
					errors["ERROR_EMAIL"] = "Email must be a valid email address"
				}
			case "Username":
				errors["ERROR_USERNAME"] = "Username is required"
			case "Name":
				errors["ERROR_NAME"] = "Name is required"
			case "Password":
				if err.Tag() == "required" {
					errors["ERROR_PASSWORD"] = "Password is required"
				} else if err.Tag() == "min" {
					errors["ERROR_PASSWORD"] = "Password must be at least 8 characters!"
				}
			case "Gender":
				errors["ERROR_GENDER"] = "Gender is required"
			}
		}

		// Store errors in cookies
		for key, msg := range errors {
			session.Set(key, msg)
			err := session.Save()
			if err != nil {
				session.Set("ERROR", err.Error())
				session.Save()
				c.Redirect(http.StatusFound, "/administrator/add-new-user")
			}
		}

		// Redirect once after storing errors
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// HASH PASSWORD
	if err := user.HashPassword(); err != nil {
		session.Set("ERROR", "Failed to hash password")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// DUPLICATE USER
	var existingUser models.User
	if err := config.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err == nil {
		if existingUser.Email == user.Email {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
		}
		if existingUser.Username == user.Username {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Activate the user
	user.IsActive = 1

	// Save user to database
	if err := config.DB.Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", "Failed to create user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_REGISTER", "New user successfully created")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func EditUser(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)

	var users []models.User
	config.DB.Where("is_active = ? AND id = ?", 1, DecodeID).Find(&users)

	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "show_user.html", gin.H{
		"title":     "Detail User",
		"menu":      menu,
		"menus":     menus,
		"user":      userSession,
		"csrfToken": csrf.GetToken(c),
		"DataUsers": users,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)
	var user models.User
	config.DB.First(&user, DecodeID)

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&user)
	c.Redirect(http.StatusFound, "/")
}

func ShowUser(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)

	var user models.User
	config.DB.First(&user, DecodeID)

	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "show_user.html", gin.H{
		"title":     "Detail User",
		"menu":      menu,
		"menus":     menus,
		"user":      userSession,
		"csrfToken": csrf.GetToken(c),
		"DataUser":  user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)
	config.DB.Unscoped().Delete(&models.User{}, DecodeID)
	c.SetCookie("SUCCESS_DELETE", "User successfully deleted", 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func ManageRole(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "manage_role.html", gin.H{
		"title": "Manage Role",
		"menu":  menu,
		"menus": menus,
		"user":  userSession,
	})
}

func ManageMenu(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "manage_menu.html", gin.H{
		"title": "Manage Menu",
		"menu":  menu,
		"menus": menus,
		"user":  userSession,
	})
}

func ManageSubmenu(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "manage_submenu.html", gin.H{
		"title": "Manage Submenu",
		"menu":  menu,
		"menus": menus,
		"user":  userSession,
	})
}
