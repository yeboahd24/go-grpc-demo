package main

import (
	"context"
	"log"
	"time"

	pb "grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to Order Service
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	orderClient := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Try to create an order
	orderReq := &pb.OrderRequest{
		ProductId: "P1",
		Quantity:  2,
		UserId:    "USER123",
	}

	response, err := orderClient.CreateOrder(ctx, orderReq)
	if err != nil {
		log.Fatalf("Could not create order: %v", err)
	}

	log.Printf("Order Response: ID=%s, Status=%s, Total Price=%.2f",
		response.OrderId, response.Status, response.TotalPrice)
}
