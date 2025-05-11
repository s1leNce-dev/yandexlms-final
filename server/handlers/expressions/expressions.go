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

func GetTask(c *gin.Context, db *gorm.DB) {
	var task models.Task
	if err := db.Where("done = ?", false).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks available"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

func SubmitTaskResult(c *gin.Context, db *gorm.DB) {
	var res struct {
		ID     uint    `json:"id"`
		Result float64 `json:"result"`
	}
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid data"})
		return
	}

	var task models.Task
	if err := db.First(&task, res.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Done = true
	db.Save(&task)

	var expr models.Expression
	db.Preload("Tasks").First(&expr, task.ExpressionID)

	allDone := true
	for _, t := range expr.Tasks {
		if !t.Done {
			allDone = false
			break
		}
	}

	expr.Result = res.Result
	if allDone {
		expr.Status = "done"
	}
	db.Save(&expr)

	c.JSON(http.StatusOK, gin.H{"message": "Result saved"})
}
