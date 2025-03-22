package helpers

import (
	"github.com/gin-gonic/gin"
	"strings"
	"unicode"
)

func Proper(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func GetMenuSubmenu(c *gin.Context) (menu, submenu string) {
	URL := strings.Split(c.Request.URL.Path, "/")
	menu = Proper(URL[1])
	submenu = Proper(URL[2])
	return menu, submenu
}
