package main

import (
	pb "awesomeProject3/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "grpc-server:50051" // Имя сервиса в Docker Compose и порт gRPC сервера
)

// main функция запускает клиент gRPC агента
func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	log.Println("Agent started")

	for {
		task, err := getTask(client)
		if err != nil {
			log.Printf("Could not get task: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		result := performTask(task)
		err = postResult(client, task.Id, result)
		if err != nil {
			log.Printf("Could not post result: %v", err)
		}
	}
}

// getTask запрашивает задачу у сервера gRPC
func getTask(client pb.TaskServiceClient) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.GetTask(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// performTask выполняет вычисление задачи
func performTask(task *pb.Task) float64 {
	// Симуляция времени выполнения задачи
	time.Sleep(time.Duration(task.OpTime) * time.Millisecond)

	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		result = task.Arg1 / task.Arg2
	}
	return result
}

// postResult отправляет результат выполнения задачи на сервер gRPC
func postResult(client pb.TaskServiceClient, id uint32, result float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.PostResult(ctx, &pb.Result{Id: id, Result: result})
	return err
}
