package helpers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/speps/go-hashids"
	"golang-example-app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type SidebarSubMenu struct {
	Id   int
	Name string
	URL  string
	Icon string
}

type SidebarMenu struct {
	Id       int
	Name     string
	SubMenus []SidebarSubMenu
}

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

func GetMenusByRole(db *gorm.DB, roleID int) []models.Menu {
	var menus []models.Menu

	tx := db.
		Table("menus").
		Joins("JOIN user_access_menus uam ON uam.menu_id = menus.id").
		Where("uam.role_id = ? AND menus.deleted_at IS NULL", roleID).
		Select("menus.id, menus.name").
		Find(&menus)

	if tx.Error != nil {
		log.Println("GetMenusByRole error:", tx.Error)
	}
	log.Println("GetMenusByRole rows:", tx.RowsAffected)

	return menus
}

func GetSubmenusByRole(db *gorm.DB, roleID int) []models.SubMenu {
	var submenus []models.SubMenu

	tx := db.
		Table("sub_menus").
		Joins("JOIN user_access_submenus uas ON uas.submenu_id = sub_menus.id").
		Where("uas.role_id = ? AND sub_menus.deleted_at IS NULL AND sub_menus.is_active = 1", roleID).
		Select("sub_menus.id, sub_menus.name, sub_menus.url, sub_menus.icon, sub_menus.menu_id").
		Find(&submenus)

	if tx.Error != nil {
		log.Println("GetSubmenusByRole error:", tx.Error)
	}
	log.Println("GetSubmenusByRole rows:", tx.RowsAffected)

	return submenus
}

func GetSidebarMenusByRole(db *gorm.DB, roleID int) []SidebarMenu {
	menus := GetMenusByRole(db, roleID)
	submenus := GetSubmenusByRole(db, roleID)

	// Map submenus by menu_id
	subMap := make(map[int][]SidebarSubMenu)
	for _, sm := range submenus {
		subMap[sm.MenuID] = append(subMap[sm.MenuID], SidebarSubMenu{
			Name: sm.Name,
			URL:  sm.URL,
			Icon: sm.Icon,
		})
	}

	// Combine into SidebarMenu
	var sidebarMenus []SidebarMenu
	for _, m := range menus {
		sidebarMenus = append(sidebarMenus, SidebarMenu{
			Name:     m.Name,
			SubMenus: subMap[m.Id],
		})
	}

	return sidebarMenus
}

func FlashMessage(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	val := session.Get(key)
	session.Delete(key)
	session.Save()
	return val
}

func GetSessionValue(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	val := session.Get(key)
	return val
}

func EncodeID(id int) string {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("SALT")
	hd.MinLength = 6 // Optional: Min length of encoded string
	h, _ := hashids.NewWithData(hd)

	e, _ := h.Encode([]int{id})
	return e
}

func DecodeID(encoded string) (int, error) {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("SALT")
	hd.MinLength = 6
	h, _ := hashids.NewWithData(hd)

	ids, err := h.DecodeWithError(encoded)
	if err != nil || len(ids) == 0 {
		return 0, fmt.Errorf("invalid ID")
	}

	return ids[0], nil
}

func GetFirstAccessibleURL(db *gorm.DB, roleID int) string {
	var submenu models.SubMenu
	err := db.
		Table("sub_menus").
		Select("sub_menus.url").
		Joins("JOIN user_access_submenus uas ON uas.submenu_id = sub_menus.id").
		Where("uas.role_id = ? AND sub_menus.is_active = 1", roleID).
		Order("sub_menus.id ASC").
		Limit(1).
		Scan(&submenu).Error

	if err != nil || submenu.URL == "" {
		return "/user/logout" // fallback jika tidak ada akses
	}

	return submenu.URL
}

func RedirectSlashRoute(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusFound, url)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
