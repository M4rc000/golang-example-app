package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"net/http"
)

func Employees(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "employees.html", gin.H{
		"title": "Employees",
		"menu":  menu,
		"menus": menus,
		"user":  userSession,
	})
}

func Attendance(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menu, _ := helpers.GetMenuSubmenu(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusOK, "attendance.html", gin.H{
		"title": "Attendance",
		"menu":  menu,
		"menus": menus,
		"user":  userSession,
	})
}
