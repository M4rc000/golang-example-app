package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang-example-app/config"
	"golang-example-app/models"
	"net/http"
)

func Dashboard(c *gin.Context) {
	session := sessions.Default(c)
	userEmailUsernameInterface := session.Get("USER_EMAIL_USERNAME")

	// Type assertion untuk memastikan userEmailUsername adalah string
	userEmailUsername, ok := userEmailUsernameInterface.(string)
	if !ok {
		// Handle error jika type assertion gagal
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Unauthorized Access"})
		c.Abort()
		return
	}

	var user models.User
	if err := config.DB.Where("email = ? OR username = ?", userEmailUsername, userEmailUsername).First(&user).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		c.Abort()
		return
	}

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

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Dashboard",
		"user":  user,
	})
}
