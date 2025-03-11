package controllers

import (
	"golang-example-app/config"
	"golang-example-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)

	// c.JSON(http.StatusOK, users)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "User List",
		"users": users,
	})
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&user)
	c.Redirect(http.StatusFound, "/users")
}

func EditUserForm(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	config.DB.First(&user, id)
	c.HTML(http.StatusOK, "edit.html", gin.H{"user": user})
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
	c.Redirect(http.StatusFound, "/users")
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.User{}, id)
	c.Redirect(http.StatusFound, "/users")
}
