package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/models"
	"net/http"
)

func AuthRequired(c *gin.Context) {
	UserID := helpers.GetSessionValue(c, "USER_ID")
	UserEmail := helpers.GetSessionValue(c, "USER_EMAIL")
	UserUsername := helpers.GetSessionValue(c, "USER_USERNAME")

	if UserID == "" || UserEmail == "" || UserUsername == "" {
		c.Redirect(http.StatusFound, "/auth")
		return
	}
	c.Next()
}

func GuestRequired(c *gin.Context) {
	UserID := helpers.GetSessionValue(c, "USER_ID")
	UserEmail := helpers.GetSessionValue(c, "USER_EMAIL")
	UserUsername := helpers.GetSessionValue(c, "USER_USERNAME")

	if UserID != "" || UserEmail != "" || UserUsername != "" {
		c.Redirect(http.StatusFound, "/administrator")
		return
	}
	c.Next()
}

func GetSessionUser(c *gin.Context) *models.User {
	UserID := helpers.GetSessionValue(c, "USER_ID")
	UserEmail := helpers.GetSessionValue(c, "USER_EMAIL")
	UserUsername := helpers.GetSessionValue(c, "USER_USERNAME")

	var user models.User
	if err := config.DB.Where("id = ? OR username = ? OR email = ?", UserID, UserUsername, UserEmail).First(&user).Error; err != nil {
		c.Redirect(http.StatusFound, "/auth")
		return nil
	}
	return &user
}
