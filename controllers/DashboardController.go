package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-example-app/config"
	"golang-example-app/helpers"
	"golang-example-app/middlewares"
	"net/http"
)

func Dashboard(c *gin.Context) {
	userSession := middlewares.GetSessionUser(c)
	menus := helpers.GetSidebarMenusByRole(config.DB, userSession.RoleID)

	c.HTML(http.StatusFound, "dashboard.html", gin.H{
		"title": "Dashboard",
		"menus": menus,
		"user":  userSession,
	})
}
