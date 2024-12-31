package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-demo/proto"
	"google.golang.org/grpc"
)

type productServer struct {
	pb.UnimplementedProductServiceServer
	// Simulated product database
	products map[string]*Product
}

type Product struct {
	ID       string
	Name     string
	Price    float32
	Quantity int32
}

func newProductServer() *productServer {
	// Initialize with some dummy data
	return &productServer{
		products: map[string]*Product{
			"P1": {ID: "P1", Name: "Laptop", Price: 999.99, Quantity: 10},
			"P2": {ID: "P2", Name: "Phone", Price: 599.99, Quantity: 20},
		},
	}
}

func (s *productServer) CheckProductAvailability(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	log.Printf("Received product availability check for ID: %v, Quantity: %v", req.ProductId, req.Quantity)
	
	product, exists := s.products[req.ProductId]
	if !exists {
		return &pb.ProductResponse{IsAvailable: false}, nil
	}

	isAvailable := product.Quantity >= req.Quantity
	return &pb.ProductResponse{
		IsAvailable:  isAvailable,
		Price:        product.Price,
		ProductName:  product.Name,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, newProductServer())
	
	fmt.Println("Product Service is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
