package grpcservice

import (
	"context"
	"log"
	"net"
	"testing"

	"final/models"
	expressionpb "final/proto/expression"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	go func() {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to open db: %v", err)
		}
		if err := db.AutoMigrate(&models.Expression{}, &models.Task{}); err != nil {
			log.Fatalf("migrate: %v", err)
		}
		grpcServer := grpc.NewServer()
		expressionpb.RegisterTaskServiceServer(grpcServer, NewServer(db))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestIntegration_GetAndSubmit(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := expressionpb.NewTaskServiceClient(conn)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&models.Expression{}, &models.Task{})
	expr := models.Expression{Status: "pending"}
	db.Create(&expr)
	task := models.Task{ExpressionID: expr.ID, Expression: "3 * 3", OperationTime: 0}
	db.Create(&task)

	resp, err := client.GetTask(ctx, &expressionpb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, uint32(task.ID), resp.Task.Id)
	assert.Equal(t, "3 * 3", resp.Task.Expression)

	_, err = client.SubmitTaskResult(ctx, &expressionpb.SubmitTaskRequest{
		Id:     resp.Task.Id,
		Result: 9,
	})
	assert.NoError(t, err)

	var updatedExpr models.Expression
	db.First(&updatedExpr, expr.ID)
	assert.Equal(t, "done", updatedExpr.Status)
	assert.Equal(t, 9.0, updatedExpr.Result)
}
