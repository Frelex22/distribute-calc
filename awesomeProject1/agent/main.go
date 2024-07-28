package main

import (
	pb "awesomeProject3/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	for {
		task, err := getTask(client)
		if err != nil {
			log.Printf("Could not get task: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		result := performTask(task)
		postResult(client, task.Id, result)
	}
}

func getTask(client pb.TaskServiceClient) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.GetTask(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func performTask(task *pb.Task) float64 {
	// Simulate task execution time
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

func postResult(client pb.TaskServiceClient, id uint32, result float64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.PostResult(ctx, &pb.Result{Id: id, Result: result})
	if err != nil {
		log.Printf("Could not post result: %v", err)
	}
}
