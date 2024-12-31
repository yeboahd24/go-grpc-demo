package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type orderServer struct {
	pb.UnimplementedOrderServiceServer
	productClient pb.ProductServiceClient
}

// gRPC service method
func (s *orderServer) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	log.Printf("Received gRPC order request for ProductID: %v, Quantity: %v", req.ProductId, req.Quantity)
	return s.processOrder(ctx, req)
}

// Common order processing logic
func (s *orderServer) processOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	// Check product availability with the Product Service
	productReq := &pb.ProductRequest{
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	}

	productResp, err := s.productClient.CheckProductAvailability(ctx, productReq)
	if err != nil {
		log.Printf("Error checking product availability: %v", err)
		return nil, err
	}

	if !productResp.IsAvailable {
		return &pb.OrderResponse{
			Status:     "FAILED",
			OrderId:    "",
			TotalPrice: 0,
		}, nil
	}

	// If product is available, create the order
	totalPrice := productResp.Price * float32(req.Quantity)
	orderId := fmt.Sprintf("ORD-%s-%s", req.UserId, req.ProductId)

	return &pb.OrderResponse{
		OrderId:    orderId,
		Status:     "SUCCESS",
		TotalPrice: totalPrice,
	}, nil
}

// REST handler
type RESTOrderRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

func (s *orderServer) handleRESTCreateOrder(c *gin.Context) {
	var restReq RESTOrderRequest
	if err := c.ShouldBindJSON(&restReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert REST request to gRPC request
	grpcReq := &pb.OrderRequest{
		ProductId: restReq.ProductID,
		Quantity:  restReq.Quantity,
		UserId:    restReq.UserID,
	}

	// Process the order using common logic
	response, err := s.processOrder(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return REST response
	c.JSON(http.StatusOK, gin.H{
		"order_id":     response.OrderId,
		"status":       response.Status,
		"total_price":  response.TotalPrice,
	})
}

func main() {
	// Connect to Product Service
	productConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Product Service: %v", err)
	}
	defer productConn.Close()
	productClient := pb.NewProductServiceClient(productConn)

	// Create order server instance
	orderSrv := &orderServer{
		productClient: productClient,
	}

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterOrderServiceServer(s, orderSrv)

		fmt.Println("gRPC Order Service is running on :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Setup REST server with Gin
	router := gin.Default()
	router.POST("/api/orders", orderSrv.handleRESTCreateOrder)

	// Start REST server
	fmt.Println("REST API Server is running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run REST server: %v", err)
	}
}
