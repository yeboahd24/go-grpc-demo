# Hybrid gRPC and REST Microservices Example

This project demonstrates a hybrid microservices architecture using both gRPC (for internal service communication) and REST APIs (for external clients) in Go. It includes a Product Service and an Order Service that communicate with each other via gRPC, while the Order Service also exposes REST endpoints for external clients.

## Project Structure

```
.
├── proto/
│   ├── product.proto        # Product service protocol buffer definition
│   └── order.proto         # Order service protocol buffer definition
├── product-service/
│   └── main.go            # Product service implementation
├── order-service/
│   └── main.go            # Order service implementation (gRPC + REST)
└── test-client/
    └── main.go            # gRPC test client
```

## Prerequisites

- Go 1.16 or later
- Protocol Buffers compiler (protoc)
- Go gRPC tools

### Installing Prerequisites

1. Install Protocol Buffers compiler:
```bash
apt-get install -y protobuf-compiler
```

2. Install Go gRPC tools:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd grpc
```

2. Generate Protocol Buffer code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/product.proto proto/order.proto
```

3. Install dependencies:
```bash
go mod tidy
```

## Running the Services

1. Start the Product Service:
```bash
go run product-service/main.go
```
This will start the Product Service on port 50051.

2. Start the Order Service (in a new terminal):
```bash
go run order-service/main.go
```
This will start:
- gRPC server on port 50052
- REST server on port 8080

## Testing the Services

### 1. Testing gRPC Endpoint

Run the test client:
```bash
go run test-client/main.go
```

Expected output:
```
2024/12/31 08:00:17 Order Response: ID=ORD-USER123-P1, Status=SUCCESS, Total Price=1999.98
```

### 2. Testing REST Endpoint

Use curl to test the REST API:
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"product_id":"P1","quantity":2,"user_id":"USER123"}'
```

Expected output:
```json
{
  "order_id": "ORD-USER123-P1",
  "status": "SUCCESS",
  "total_price": 1999.98
}
```

## API Documentation

### gRPC Services

#### Product Service
- Endpoint: localhost:50051
- Methods:
  - `CheckProductAvailability(ProductRequest) returns (ProductResponse)`

#### Order Service (gRPC)
- Endpoint: localhost:50052
- Methods:
  - `CreateOrder(OrderRequest) returns (OrderResponse)`

### REST API

#### Create Order
- Endpoint: POST /api/orders
- Request Body:
```json
{
  "product_id": "string",
  "quantity": number,
  "user_id": "string"
}
```
- Response:
```json
{
  "order_id": "string",
  "status": "string",
  "total_price": number
}
```

## When to Use This Architecture

This hybrid architecture is particularly suitable for:

1. **Large-Scale Microservices Applications**
   - Internal Services: Use gRPC for high-performance communication
   - External API Gateway: REST APIs for third-party integrations

2. **Performance-Critical Applications**
   - Fast internal communication (gRPC)
   - Broad client compatibility (REST)

3. **Legacy System Integration**
   - New microservices: gRPC
   - Legacy clients: REST

4. **Public API Providers**
   - Internal efficiency with gRPC
   - Developer-friendly REST APIs

5. **Mobile/Web Applications with Complex Backend**
   - Backend Services: gRPC
   - Client Applications: REST

### Key Benefits

1. **Performance**
   - Binary protocol with Protocol Buffers
   - HTTP/2 streaming
   - Efficient inter-service communication

2. **Developer Experience**
   - Strong typing with gRPC
   - Familiar REST APIs for external developers

3. **Flexibility**
   - Choose right protocol for each use case
   - Easy integration with different clients

4. **Scalability**
   - Efficient high-volume internal traffic
   - Suitable for various external clients

5. **Maintainability**
   - Clear separation of internal/external APIs
   - Shared business logic
   - Type safety for internal communications

### When Not to Use

1. Simple applications or monoliths
2. Small teams with limited resources
3. Uniform client base that can all use gRPC
4. Resource-constrained environments

## Example Use Case

```plaintext
E-commerce Platform:
├── Internal Services (gRPC)
│   ├── Order Service
│   ├── Inventory Service
│   ├── Payment Service
│   └── User Service
│
└── External APIs (REST)
    ├── Mobile App API
    ├── Web Client API
    ├── Partner Integration API
    └── Admin Dashboard API
```

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
