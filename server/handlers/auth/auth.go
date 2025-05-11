package auth

import (
	"final/models"
	"final/utils/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	secure := gin.Mode() != gin.DebugMode
	sameSite := http.SameSiteNoneMode
	if !secure {
		sameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_access",
		Value:    accessToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_refresh",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}

func Login(c *gin.Context, db *gorm.DB) {
	var loginform struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginform); err != nil {
		c.JSON(400, gin.H{"error": "eof"})
		return
	}

	var user models.User
	result := db.Where("login = ?", loginform.Login).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": "record not found"})
		return
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Login)
	if err != nil {
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	setAuthCookies(c, accessToken, refreshToken)

	c.JSON(200, gin.H{"message": "success"})
}

func Regin(c *gin.Context, db *gorm.DB) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(500, gin.H{"error": "EOF"})
		return
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to create"})
		return
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(newUser.ID, newUser.Login)
	if err != nil {
		c.JSON(500, gin.H{"error": "token error"})
		db.Unscoped().Where("id = ?", newUser.ID).Delete(&models.User{})
		return
	}

	setAuthCookies(c, accessToken, refreshToken)

	c.JSON(200, gin.H{"message": "success"})
}
