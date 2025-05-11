package authmiddleware

import (
	jwtutils "final/utils/jwt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_access")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "empty token",
			})
			c.Abort()
			return
		}

		jwtClaims, err := jwtutils.ValidateAccessToken(tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token signature",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "token error",
			})
			log.Println("[ERROR]", err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", jwtClaims.UserID)
		c.Set("login", jwtClaims.Login)
		c.Next()
	}
}
