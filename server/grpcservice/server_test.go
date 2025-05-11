package grpcservice

import (
	"context"
	"testing"

	"final/models"
	expressionpb "final/proto/expression"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Expression{}, &models.Task{})
	assert.NoError(t, err)
	return db
}

func TestGetTask_NoTasks(t *testing.T) {
	db := setupInMemoryDB(t)
	srv := NewServer(db)

	_, err := srv.GetTask(context.Background(), &expressionpb.Empty{})
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGetTask_WithPendingTask(t *testing.T) {
	db := setupInMemoryDB(t)

	exp := models.Expression{Status: "pending"}
	db.Create(&exp)
	task := models.Task{
		ExpressionID:  exp.ID,
		Expression:    "1 + 2",
		OperationTime: 100,
		Done:          false,
	}
	db.Create(&task)

	srv := NewServer(db)
	resp, err := srv.GetTask(context.Background(), &expressionpb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, uint32(task.ID), resp.Task.Id)
	assert.Equal(t, "1 + 2", resp.Task.Expression)
	assert.Equal(t, int32(100), resp.Task.OperationTime)
}

func TestSubmitTaskResult_MarkDoneAndFinishExpression(t *testing.T) {
	db := setupInMemoryDB(t)

	exp := models.Expression{Status: "pending"}
	db.Create(&exp)
	t1 := models.Task{ExpressionID: exp.ID, Expression: "1 + 1", OperationTime: 0}
	t2 := models.Task{ExpressionID: exp.ID, Expression: "2 + 2", OperationTime: 0}
	db.Create(&t1)
	db.Create(&t2)

	srv := NewServer(db)

	_, err := srv.SubmitTaskResult(context.Background(), &expressionpb.SubmitTaskRequest{
		Id:     uint32(t1.ID),
		Result: 2,
	})
	assert.NoError(t, err)

	var expr models.Expression
	db.First(&expr, exp.ID)
	assert.Equal(t, "pending", expr.Status)

	_, err = srv.SubmitTaskResult(context.Background(), &expressionpb.SubmitTaskRequest{
		Id:     uint32(t2.ID),
		Result: 4,
	})
	assert.NoError(t, err)

	db.First(&expr, exp.ID)
	assert.Equal(t, "done", expr.Status)
	assert.Equal(t, 4.0, expr.Result)
}
