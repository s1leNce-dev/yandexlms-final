package expressions

import (
	"final/models"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func parseExpression(exprStr string) []models.Task {
	tokens := strings.Fields(exprStr)
	if len(tokens) < 3 {
		return nil
	}

	var tasks []models.Task
	for i := 1; i < len(tokens)-1; i += 2 {
		exprPart := tokens[i-1] + " " + tokens[i] + " " + tokens[i+1]
		tasks = append(tasks, models.Task{
			Expression:    exprPart,
			OperationTime: rand.Intn(5000) + 1000,
		})
	}
	return tasks
}

func AddExpression(c *gin.Context, db *gorm.DB) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid expression"})
		return
	}

	tasks := parseExpression(req.Expression)
	if tasks == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid expression format"})
		return
	}

	expr := models.Expression{
		Status: "pending",
		Tasks:  tasks,
	}
	db.Create(&expr)

	c.JSON(http.StatusCreated, gin.H{"id": expr.ID})
}

func GetExpressions(c *gin.Context, db *gorm.DB) {
	var expressions []models.Expression
	db.Find(&expressions)
	c.JSON(http.StatusOK, gin.H{"expressions": expressions})
}

func GetExpressionByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var expr models.Expression
	if err := db.Preload("Tasks").First(&expr, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expression not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expression": expr})
}
