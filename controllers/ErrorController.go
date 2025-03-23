package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoFoundRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "NotFound.html", gin.H{
		"title": "Not Found 404",
	})
}
