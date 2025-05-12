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
	successUpdate := helpers.FlashMessage(c, "SUCCESS_UPDATE")
	successDelete := helpers.FlashMessage(c, "SUCCESS_DELETE")

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
		"successUpdate":   successUpdate,
		"successDelete":   successDelete,
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
	oldName := helpers.FlashMessage(c, "OLD_NAME")
	oldUsername := helpers.FlashMessage(c, "OLD_USERNAME")
	oldEmail := helpers.FlashMessage(c, "OLD_EMAIL")
	oldGender := helpers.FlashMessage(c, "OLD_GENDER")

	c.HTML(http.StatusOK, "add_user.html", gin.H{
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
		"oldName":           oldName,
		"oldUsername":       oldUsername,
		"oldEmail":          oldEmail,
		"oldGender":         oldGender,
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

	// OLD VALUE
	session.Set("OLD_NAME", user.Name)
	session.Set("OLD_USERNAME", user.Username)
	session.Set("OLD_EMAIL", user.Email)
	session.Set("OLD_GENDER", user.Gender)
	session.Save()

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
	password, err := helpers.HashPassword(user.Password)
	if err != nil {
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
	user.Password = password

	// Save user to database
	if err := config.DB.Table("users").Create(&user).Error; err != nil {
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

	var users models.User
	config.DB.Where("is_active = ? AND id = ?", 1, DecodeID).First(&users)

	EncodedID := helpers.EncodeID(users.Id)

	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)
	err := helpers.FlashMessage(c, "ERROR")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	duplicateEmail := helpers.FlashMessage(c, "DUPLICATE_EMAIL")
	duplicateUsername := helpers.FlashMessage(c, "DUPLICATE_USERNAME")
	failedUpdate := helpers.FlashMessage(c, "FAILED_UP")
	errorEmail := helpers.FlashMessage(c, "ERROR_EMAIL")
	errorUsername := helpers.FlashMessage(c, "ERROR_USERNAME")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorPassword := helpers.FlashMessage(c, "ERROR_PASSWORD")
	errorGender := helpers.FlashMessage(c, "ERROR_GENDER")

	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"title":             "Edit User",
		"menu":              menu,
		"menus":             menus,
		"user":              userSession,
		"csrfToken":         csrf.GetToken(c),
		"users":             users,
		"EncodeID":          EncodedID,
		"err":               err,
		"errorInputData":    errorInputData,
		"duplicateUsername": duplicateUsername,
		"duplicateEmail":    duplicateEmail,
		"failedUpdate":      failedUpdate,
		"errorEmail":        errorEmail,
		"errorUsername":     errorUsername,
		"errorName":         errorName,
		"errorPassword":     errorPassword,
		"errorGender":       errorGender,
	})
}

func UpdateUser(c *gin.Context) {
	session := sessions.Default(c)
	validate := validator.New()
	errors := make(map[string]string)

	// Get and decode ID
	id := c.Param("id")
	DecodeID, err := helpers.DecodeID(id)
	if err != nil {
		session.Set("ERROR", "Invalid user ID")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Fetch existing user
	var user models.User
	if err := config.DB.First(&user, DecodeID).Error; err != nil {
		session.Set("ERROR", "User not found")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Bind input to temporary struct
	var input models.UpdateManageUser
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/edit-new-user/"+id)
		return
	}

	// Validate input
	if err := validate.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errors["ERROR_EMAIL"] = "Email must be a valid email address"
			case "Username":
				errors["ERROR_USERNAME"] = "Username is required"
			case "Name":
				errors["ERROR_NAME"] = "Name is required"
			case "Password":
				if err.Tag() == "required" {
					errors["ERROR_PASSWORD"] = "Password is required"
				} else if err.Tag() == "min" {
					errors["ERROR_PASSWORD"] = "Password must be at least 8 characters"
				}
			case "Gender":
				errors["ERROR_GENDER"] = "Gender is required"
			}
		}

		// Store validation errors in session
		for key, msg := range errors {
			session.Set(key, msg)
			session.Save()
		}

		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Check for duplicate email or username (excluding current user)
	var existingUser models.User
	if err := config.DB.Where("(email = ? OR username = ?) AND id != ?", input.Email, input.Username, DecodeID).First(&existingUser).Error; err == nil {
		if existingUser.Email == input.Email {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
		}
		if existingUser.Username == input.Username {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Hash the password only if it's not empty
	if input.Password != "" {
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			session.Set("ERROR", "Failed to hash password")
			session.Save()
			c.Redirect(http.StatusFound, "/administrator/add-new-user")
			return
		}
		input.Password = hashedPassword
	} else {
		input.Password = user.Password // Keep old password if not changed
	}

	// Update fields
	user.Email = input.Email
	user.Username = input.Username
	user.Name = input.Name
	user.Password = input.Password
	user.Gender = input.Gender

	if err := config.DB.Table("users").Save(&user).Error; err != nil {
		session.Set("FAILED_UPDATE", "Failed to update user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_UPDATE", "User successfully updated")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
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
	session := sessions.Default(c)
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)
	config.DB.Unscoped().Delete(&models.User{}, DecodeID)

	session.Set("SUCCESS_DELETE", "User successfully deleted")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func ManageRole(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	err := helpers.FlashMessage(c, "ERROR")
	successRegister := helpers.FlashMessage(c, "SUCCESS_REGISTER")
	successUpdate := helpers.FlashMessage(c, "SUCCESS_UPDATE")
	successDelete := helpers.FlashMessage(c, "SUCCESS_DELETE")

	var roles []models.Role
	config.DB.Find(&roles)

	var DataUsers []map[string]interface{}
	for i, role := range roles {
		DataUsers = append(DataUsers, map[string]interface{}{
			"Number":    i + 1,
			"Id":        role.Id,
			"EncodedId": helpers.EncodeID(role.Id),
			"Name":      role.Name,
			"Active":    role.IsActive,
		})
	}

	c.HTML(http.StatusOK, "manage_role.html", gin.H{
		"title":           "Manage Role",
		"menu":            menu,
		"menus":           menus,
		"user":            userSession,
		"err":             err,
		"successRegister": successRegister,
		"successUpdate":   successUpdate,
		"successDelete":   successDelete,
		"DataUsers":       DataUsers,
	})
}

func AddNewRole(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	err := helpers.FlashMessage(c, "ERROR")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	duplicateName := helpers.FlashMessage(c, "DUPLICATE_NAME")
	failedRegister := helpers.FlashMessage(c, "FAILED_REGISTER")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorIsActive := helpers.FlashMessage(c, "ERROR_ISACTIVE")
	oldName := helpers.FlashMessage(c, "OLD_NAME")
	oldIsActive := helpers.FlashMessage(c, "OLD_ISACTIVE")

	c.HTML(http.StatusOK, "add_role.html", gin.H{
		"title":          "New Role",
		"menu":           menu,
		"menus":          menus,
		"user":           userSession,
		"csrfToken":      csrf.GetToken(c),
		"err":            err,
		"errorInputData": errorInputData,
		"duplicateName":  duplicateName,
		"failedRegister": failedRegister,
		"errorIsActive":  errorIsActive,
		"errorName":      errorName,
		"oldName":        oldName,
		"oldIsActive":    oldIsActive,
	})
}

func SaveNewRole(c *gin.Context) {
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

	// OLD VALUE
	session.Set("OLD_NAME", user.Name)
	session.Set("OLD_USERNAME", user.Username)
	session.Set("OLD_EMAIL", user.Email)
	session.Set("OLD_GENDER", user.Gender)
	session.Save()

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
	password, err := helpers.HashPassword(user.Password)
	if err != nil {
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
	user.Password = password

	// Save user to database
	if err := config.DB.Table("users").Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", "Failed to create user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_REGISTER", "New user successfully created")
	session.Save()

	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func EditRole(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)

	var users models.User
	config.DB.Where("is_active = ? AND id = ?", 1, DecodeID).First(&users)

	EncodedID := helpers.EncodeID(users.Id)

	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)
	err := helpers.FlashMessage(c, "ERROR")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	duplicateEmail := helpers.FlashMessage(c, "DUPLICATE_EMAIL")
	duplicateUsername := helpers.FlashMessage(c, "DUPLICATE_USERNAME")
	failedUpdate := helpers.FlashMessage(c, "FAILED_UP")
	errorEmail := helpers.FlashMessage(c, "ERROR_EMAIL")
	errorUsername := helpers.FlashMessage(c, "ERROR_USERNAME")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorPassword := helpers.FlashMessage(c, "ERROR_PASSWORD")
	errorGender := helpers.FlashMessage(c, "ERROR_GENDER")

	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"title":             "Edit User",
		"menu":              menu,
		"menus":             menus,
		"user":              userSession,
		"csrfToken":         csrf.GetToken(c),
		"users":             users,
		"EncodeID":          EncodedID,
		"err":               err,
		"errorInputData":    errorInputData,
		"duplicateUsername": duplicateUsername,
		"duplicateEmail":    duplicateEmail,
		"failedUpdate":      failedUpdate,
		"errorEmail":        errorEmail,
		"errorUsername":     errorUsername,
		"errorName":         errorName,
		"errorPassword":     errorPassword,
		"errorGender":       errorGender,
	})
}

func UpdateRole(c *gin.Context) {
	session := sessions.Default(c)
	validate := validator.New()
	errors := make(map[string]string)

	// Get and decode ID
	id := c.Param("id")
	DecodeID, err := helpers.DecodeID(id)
	if err != nil {
		session.Set("ERROR", "Invalid user ID")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Fetch existing user
	var user models.User
	if err := config.DB.First(&user, DecodeID).Error; err != nil {
		session.Set("ERROR", "User not found")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Bind input to temporary struct
	var input models.UpdateManageUser
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/edit-new-user/"+id)
		return
	}

	// Validate input
	if err := validate.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errors["ERROR_EMAIL"] = "Email must be a valid email address"
			case "Username":
				errors["ERROR_USERNAME"] = "Username is required"
			case "Name":
				errors["ERROR_NAME"] = "Name is required"
			case "Password":
				if err.Tag() == "required" {
					errors["ERROR_PASSWORD"] = "Password is required"
				} else if err.Tag() == "min" {
					errors["ERROR_PASSWORD"] = "Password must be at least 8 characters"
				}
			case "Gender":
				errors["ERROR_GENDER"] = "Gender is required"
			}
		}

		// Store validation errors in session
		for key, msg := range errors {
			session.Set(key, msg)
			session.Save()
		}

		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Check for duplicate email or username (excluding current user)
	var existingUser models.User
	if err := config.DB.Where("(email = ? OR username = ?) AND id != ?", input.Email, input.Username, DecodeID).First(&existingUser).Error; err == nil {
		if existingUser.Email == input.Email {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
		}
		if existingUser.Username == input.Username {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Hash the password only if it's not empty
	if input.Password != "" {
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			session.Set("ERROR", "Failed to hash password")
			session.Save()
			c.Redirect(http.StatusFound, "/administrator/add-new-user")
			return
		}
		input.Password = hashedPassword
	} else {
		input.Password = user.Password // Keep old password if not changed
	}

	// Update fields
	user.Email = input.Email
	user.Username = input.Username
	user.Name = input.Name
	user.Password = input.Password
	user.Gender = input.Gender

	if err := config.DB.Table("users").Save(&user).Error; err != nil {
		session.Set("FAILED_UPDATE", "Failed to update user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_UPDATE", "User successfully updated")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func ShowRole(c *gin.Context) {
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

func DeleteRole(c *gin.Context) {
	session := sessions.Default(c)
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)
	config.DB.Unscoped().Delete(&models.User{}, DecodeID)

	session.Set("SUCCESS_DELETE", "User successfully deleted")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
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

func AddNewMenu(c *gin.Context) {
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
	oldName := helpers.FlashMessage(c, "OLD_NAME")
	oldUsername := helpers.FlashMessage(c, "OLD_USERNAME")
	oldEmail := helpers.FlashMessage(c, "OLD_EMAIL")
	oldGender := helpers.FlashMessage(c, "OLD_GENDER")

	c.HTML(http.StatusOK, "add_user.html", gin.H{
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
		"oldName":           oldName,
		"oldUsername":       oldUsername,
		"oldEmail":          oldEmail,
		"oldGender":         oldGender,
	})
}

func SaveNewMenu(c *gin.Context) {
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

	// OLD VALUE
	session.Set("OLD_NAME", user.Name)
	session.Set("OLD_USERNAME", user.Username)
	session.Set("OLD_EMAIL", user.Email)
	session.Set("OLD_GENDER", user.Gender)
	session.Save()

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
	password, err := helpers.HashPassword(user.Password)
	if err != nil {
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
	user.Password = password

	// Save user to database
	if err := config.DB.Table("users").Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", "Failed to create user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_REGISTER", "New user successfully created")
	session.Save()

	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func EditMenu(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)

	var users models.User
	config.DB.Where("is_active = ? AND id = ?", 1, DecodeID).First(&users)

	EncodedID := helpers.EncodeID(users.Id)

	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)
	err := helpers.FlashMessage(c, "ERROR")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	duplicateEmail := helpers.FlashMessage(c, "DUPLICATE_EMAIL")
	duplicateUsername := helpers.FlashMessage(c, "DUPLICATE_USERNAME")
	failedUpdate := helpers.FlashMessage(c, "FAILED_UP")
	errorEmail := helpers.FlashMessage(c, "ERROR_EMAIL")
	errorUsername := helpers.FlashMessage(c, "ERROR_USERNAME")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorPassword := helpers.FlashMessage(c, "ERROR_PASSWORD")
	errorGender := helpers.FlashMessage(c, "ERROR_GENDER")

	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"title":             "Edit User",
		"menu":              menu,
		"menus":             menus,
		"user":              userSession,
		"csrfToken":         csrf.GetToken(c),
		"users":             users,
		"EncodeID":          EncodedID,
		"err":               err,
		"errorInputData":    errorInputData,
		"duplicateUsername": duplicateUsername,
		"duplicateEmail":    duplicateEmail,
		"failedUpdate":      failedUpdate,
		"errorEmail":        errorEmail,
		"errorUsername":     errorUsername,
		"errorName":         errorName,
		"errorPassword":     errorPassword,
		"errorGender":       errorGender,
	})
}

func UpdateMenu(c *gin.Context) {
	session := sessions.Default(c)
	validate := validator.New()
	errors := make(map[string]string)

	// Get and decode ID
	id := c.Param("id")
	DecodeID, err := helpers.DecodeID(id)
	if err != nil {
		session.Set("ERROR", "Invalid user ID")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Fetch existing user
	var user models.User
	if err := config.DB.First(&user, DecodeID).Error; err != nil {
		session.Set("ERROR", "User not found")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Bind input to temporary struct
	var input models.UpdateManageUser
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/edit-new-user/"+id)
		return
	}

	// Validate input
	if err := validate.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errors["ERROR_EMAIL"] = "Email must be a valid email address"
			case "Username":
				errors["ERROR_USERNAME"] = "Username is required"
			case "Name":
				errors["ERROR_NAME"] = "Name is required"
			case "Password":
				if err.Tag() == "required" {
					errors["ERROR_PASSWORD"] = "Password is required"
				} else if err.Tag() == "min" {
					errors["ERROR_PASSWORD"] = "Password must be at least 8 characters"
				}
			case "Gender":
				errors["ERROR_GENDER"] = "Gender is required"
			}
		}

		// Store validation errors in session
		for key, msg := range errors {
			session.Set(key, msg)
			session.Save()
		}

		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Check for duplicate email or username (excluding current user)
	var existingUser models.User
	if err := config.DB.Where("(email = ? OR username = ?) AND id != ?", input.Email, input.Username, DecodeID).First(&existingUser).Error; err == nil {
		if existingUser.Email == input.Email {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
		}
		if existingUser.Username == input.Username {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Hash the password only if it's not empty
	if input.Password != "" {
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			session.Set("ERROR", "Failed to hash password")
			session.Save()
			c.Redirect(http.StatusFound, "/administrator/add-new-user")
			return
		}
		input.Password = hashedPassword
	} else {
		input.Password = user.Password // Keep old password if not changed
	}

	// Update fields
	user.Email = input.Email
	user.Username = input.Username
	user.Name = input.Name
	user.Password = input.Password
	user.Gender = input.Gender

	if err := config.DB.Table("users").Save(&user).Error; err != nil {
		session.Set("FAILED_UPDATE", "Failed to update user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_UPDATE", "User successfully updated")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func ShowMenu(c *gin.Context) {
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

func DeleteMenu(c *gin.Context) {
	session := sessions.Default(c)
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)
	config.DB.Unscoped().Delete(&models.User{}, DecodeID)

	session.Set("SUCCESS_DELETE", "User successfully deleted")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
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

func AddNewSubmenu(c *gin.Context) {
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
	oldName := helpers.FlashMessage(c, "OLD_NAME")
	oldUsername := helpers.FlashMessage(c, "OLD_USERNAME")
	oldEmail := helpers.FlashMessage(c, "OLD_EMAIL")
	oldGender := helpers.FlashMessage(c, "OLD_GENDER")

	c.HTML(http.StatusOK, "add_user.html", gin.H{
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
		"oldName":           oldName,
		"oldUsername":       oldUsername,
		"oldEmail":          oldEmail,
		"oldGender":         oldGender,
	})
}

func SaveNewSubmenu(c *gin.Context) {
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

	// OLD VALUE
	session.Set("OLD_NAME", user.Name)
	session.Set("OLD_USERNAME", user.Username)
	session.Set("OLD_EMAIL", user.Email)
	session.Set("OLD_GENDER", user.Gender)
	session.Save()

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
	password, err := helpers.HashPassword(user.Password)
	if err != nil {
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
	user.Password = password

	// Save user to database
	if err := config.DB.Table("users").Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", "Failed to create user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_REGISTER", "New user successfully created")
	session.Save()

	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func EditSubmenu(c *gin.Context) {
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)

	var users models.User
	config.DB.Where("is_active = ? AND id = ?", 1, DecodeID).First(&users)

	EncodedID := helpers.EncodeID(users.Id)

	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)
	err := helpers.FlashMessage(c, "ERROR")
	errorInputData := helpers.FlashMessage(c, "ERROR_INPUTDATA")
	duplicateEmail := helpers.FlashMessage(c, "DUPLICATE_EMAIL")
	duplicateUsername := helpers.FlashMessage(c, "DUPLICATE_USERNAME")
	failedUpdate := helpers.FlashMessage(c, "FAILED_UP")
	errorEmail := helpers.FlashMessage(c, "ERROR_EMAIL")
	errorUsername := helpers.FlashMessage(c, "ERROR_USERNAME")
	errorName := helpers.FlashMessage(c, "ERROR_NAME")
	errorPassword := helpers.FlashMessage(c, "ERROR_PASSWORD")
	errorGender := helpers.FlashMessage(c, "ERROR_GENDER")

	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"title":             "Edit User",
		"menu":              menu,
		"menus":             menus,
		"user":              userSession,
		"csrfToken":         csrf.GetToken(c),
		"users":             users,
		"EncodeID":          EncodedID,
		"err":               err,
		"errorInputData":    errorInputData,
		"duplicateUsername": duplicateUsername,
		"duplicateEmail":    duplicateEmail,
		"failedUpdate":      failedUpdate,
		"errorEmail":        errorEmail,
		"errorUsername":     errorUsername,
		"errorName":         errorName,
		"errorPassword":     errorPassword,
		"errorGender":       errorGender,
	})
}

func UpdateSubmenu(c *gin.Context) {
	session := sessions.Default(c)
	validate := validator.New()
	errors := make(map[string]string)

	// Get and decode ID
	id := c.Param("id")
	DecodeID, err := helpers.DecodeID(id)
	if err != nil {
		session.Set("ERROR", "Invalid user ID")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Fetch existing user
	var user models.User
	if err := config.DB.First(&user, DecodeID).Error; err != nil {
		session.Set("ERROR", "User not found")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/manage-user")
		return
	}

	// Bind input to temporary struct
	var input models.UpdateManageUser
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/edit-new-user/"+id)
		return
	}

	// Validate input
	if err := validate.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errors["ERROR_EMAIL"] = "Email must be a valid email address"
			case "Username":
				errors["ERROR_USERNAME"] = "Username is required"
			case "Name":
				errors["ERROR_NAME"] = "Name is required"
			case "Password":
				if err.Tag() == "required" {
					errors["ERROR_PASSWORD"] = "Password is required"
				} else if err.Tag() == "min" {
					errors["ERROR_PASSWORD"] = "Password must be at least 8 characters"
				}
			case "Gender":
				errors["ERROR_GENDER"] = "Gender is required"
			}
		}

		// Store validation errors in session
		for key, msg := range errors {
			session.Set(key, msg)
			session.Save()
		}

		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Check for duplicate email or username (excluding current user)
	var existingUser models.User
	if err := config.DB.Where("(email = ? OR username = ?) AND id != ?", input.Email, input.Username, DecodeID).First(&existingUser).Error; err == nil {
		if existingUser.Email == input.Email {
			session.Set("DUPLICATE_EMAIL", "Email already exists")
		}
		if existingUser.Username == input.Username {
			session.Set("DUPLICATE_USERNAME", "Username already exists")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	// Hash the password only if it's not empty
	if input.Password != "" {
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			session.Set("ERROR", "Failed to hash password")
			session.Save()
			c.Redirect(http.StatusFound, "/administrator/add-new-user")
			return
		}
		input.Password = hashedPassword
	} else {
		input.Password = user.Password // Keep old password if not changed
	}

	// Update fields
	user.Email = input.Email
	user.Username = input.Username
	user.Name = input.Name
	user.Password = input.Password
	user.Gender = input.Gender

	if err := config.DB.Table("users").Save(&user).Error; err != nil {
		session.Set("FAILED_UPDATE", "Failed to update user")
		session.Save()
		c.Redirect(http.StatusFound, "/administrator/add-new-user")
		return
	}

	session.Set("SUCCESS_UPDATE", "User successfully updated")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
}

func ShowSubmenu(c *gin.Context) {
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

func DeleteSubmenu(c *gin.Context) {
	session := sessions.Default(c)
	id := c.Param("id")
	DecodeID, _ := helpers.DecodeID(id)
	config.DB.Unscoped().Delete(&models.User{}, DecodeID)

	session.Set("SUCCESS_DELETE", "User successfully deleted")
	session.Save()
	c.Redirect(http.StatusFound, "/administrator/manage-user")
}
