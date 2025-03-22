package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"net/http"
)

func Dashboard(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, submenu := helpers.GetMenuSubmenu(c)

	//config.DB.Where("is_active = ?", 1).Find(&users)

	//var DataUsers []map[string]interface{}
	//for i, user := range users {
	//	DataUsers = append(DataUsers, map[string]interface{}{
	//		"Number":   i + 1,
	//		"Id":       user.Id,
	//		"Name":     user.Name,
	//		"Email":    user.Email,
	//		"IsActive": user.IsActive,
	//	})
	//}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":   "Dashboard",
		"user":    userSession,
		"menu":    menu,
		"submenu": submenu,
	})
}
