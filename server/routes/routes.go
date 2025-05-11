package routes

import (
	"final/handlers/auth"
	"final/handlers/expressions"
	"final/middlewares/authmiddleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	api_endpoint_v1 = "/api/v1"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// auth
	r.POST(api_endpoint_v1+"/regin", func(c *gin.Context) {
		auth.Regin(c, db)
	})
	r.POST(api_endpoint_v1+"/login", func(c *gin.Context) {
		auth.Login(c, db)
	})

	userRoutes := r.Group(api_endpoint_v1 + "/user")
	userRoutes.Use(authmiddleware.AuthMiddleware())
	{
		userRoutes.POST("/calculate", func(c *gin.Context) {
			expressions.AddExpression(c, db)
		})
		userRoutes.GET("/expressions", func(c *gin.Context) {
			expressions.GetExpressions(c, db)
		})
		userRoutes.GET("/expressions/:id", func(c *gin.Context) {
			expressions.GetExpressionByID(c, db)
		})
	}

	agentRoutes := r.Group("/internal")
	{
		agentRoutes.GET("/task", func(c *gin.Context) {
			expressions.GetTask(c, db)
		})
		agentRoutes.POST("/task", func(c *gin.Context) {
			expressions.SubmitTaskResult(c, db)
		})
	}
}
