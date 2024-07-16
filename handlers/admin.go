package handlers

import (
	"fmt"
	"ginapp/domain"
	"ginapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Adminlogin(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func AdminDashboard(c *gin.Context) {
	c.HTML(200, "dashboard.html", nil)

}

func AdminLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input domain.User
		input.Email = c.PostForm("email")

		fmt.Println("here is the email", input.Email)
		input.Password = c.PostForm("password")

		var user domain.User

		if err := db.Where("email=?", input.Email).First(&user).Error; err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid credentials"})
			return
		}

		if !user.IsAdmin {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "user is not admin"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid credentials"})
			return
		}
		token, err := utils.GenerateToken(user.Email)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "token based error"})
			return
		}
		c.SetCookie("token", token, 3600, "/", "localhost", false, true)
		fmt.Println("here is the token", token)
		c.Redirect(http.StatusFound, "/admin/dashboard")

	}
}