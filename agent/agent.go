package main

import (
	"context"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"agentim/eval"
	expressionpb "agentim/proto/expression"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultWorkers  = 2
	defaultGRPCAddr = "main_app:50051"
	requestTimeout  = 5 * time.Second
	retryInterval   = 2 * time.Second
)

func calc(task *expressionpb.Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	res, err := eval.Eval(task.Expression)
	if err != nil {
		log.Printf("Ошибка вычисления '%s': %v", task.Expression, err)
		return math.NaN()
	}
	return res
}

func worker(client expressionpb.TaskServiceClient) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		resp, err := client.GetTask(ctx, &expressionpb.Empty{})
		cancel()
		if err != nil {
			if status.Code(err) == codes.NotFound {
				time.Sleep(retryInterval)
				continue
			}
			log.Printf("gRPC GetTask error: %v", err)
			time.Sleep(retryInterval)
			continue
		}

		if resp.Task == nil || resp.Task.Id == 0 {
			time.Sleep(retryInterval)
			continue
		}

		result := calc(resp.Task)
		ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
		_, err = client.SubmitTaskResult(ctx, &expressionpb.SubmitTaskRequest{
			Id:     resp.Task.Id,
			Result: result,
		})
		cancel()
		if err != nil {
			log.Printf("gRPC SubmitTaskResult error: %v", err)
		} else {
			log.Printf("Task %d done: %f", resp.Task.Id, result)
		}
	}
}

func main() {
	workers := defaultWorkers
	if v, err := strconv.Atoi(os.Getenv("COMPUTING_POWER")); err == nil && v > 0 {
		workers = v
	}

	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = defaultGRPCAddr
	}
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC-серверу %s: %v", grpcAddr, err)
	}
	defer conn.Close()

	client := expressionpb.NewTaskServiceClient(conn)
	for i := 0; i < workers; i++ {
		go worker(client)
	}
	select {}
}
