package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang-example-app/config"
	"golang-example-app/models"
	"net/http"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("USER_ID")
	if userID == nil {
		c.Redirect(http.StatusFound, "/auth")
		c.Abort()
		return // Add return to prevent further execution
	}
	c.Next()
}

func GuestRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("USER_ID")
	if userID != nil {
		c.Redirect(http.StatusFound, "/home/dashboard")
		c.Abort()
		return // Add return to prevent further execution
	}
	c.Next()
}

func GetSessionUser(c *gin.Context) *models.User {
	session := sessions.Default(c)
	userEmailUsernameInterface := session.Get("USER_EMAIL_USERNAME")

	// Type assertion untuk memastikan userEmailUsername adalah string
	userEmailUsername, ok := userEmailUsernameInterface.(string)
	if !ok {
		// Handle error jika type assertion gagal
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Unauthorized Access"})
		c.Abort()
		return nil
	}

	var user models.User
	if err := config.DB.Where("email = ? OR username = ?", userEmailUsername, userEmailUsername).First(&user).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		c.Abort()
		return nil
	}

	return &user
}
