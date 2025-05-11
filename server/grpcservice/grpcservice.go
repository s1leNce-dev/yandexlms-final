package grpcservice

import (
	"context"
	"errors"

	"final/models"
	expressionpb "final/proto/expression"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type taskServer struct {
	expressionpb.UnimplementedTaskServiceServer
	db *gorm.DB
}

func NewServer(db *gorm.DB) *taskServer {
	return &taskServer{db: db}
}

func (s *taskServer) GetTask(ctx context.Context, _ *expressionpb.Empty) (*expressionpb.GetTaskResponse, error) {
	var task models.Task
	err := s.db.Where("done = ?", false).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &expressionpb.GetTaskResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &expressionpb.GetTaskResponse{
		Task: &expressionpb.Task{
			Id:            uint32(task.ID),
			Expression:    task.Expression,
			OperationTime: uint32(task.OperationTime),
		},
	}, nil
}

func (s *taskServer) SubmitTaskResult(ctx context.Context, req *expressionpb.SubmitTaskRequest) (*expressionpb.SubmitTaskResponse, error) {
	var task models.Task
	if err := s.db.First(&task, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found")
	}

	task.Done = true
	if err := s.db.Save(&task).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save task: %v", err)
	}

	var expr models.Expression
	if err := s.db.Preload("Tasks").First(&expr, task.ExpressionID).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load expression: %v", err)
	}

	allDone := true
	for _, t := range expr.Tasks {
		if !t.Done {
			allDone = false
			break
		}
	}

	expr.Result = req.Result
	if allDone {
		expr.Status = "done"
	}
	if err := s.db.Save(&expr).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update expression: %v", err)
	}

	return &expressionpb.SubmitTaskResponse{
		Message: "Result saved",
	}, nil
}
