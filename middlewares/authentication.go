package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
