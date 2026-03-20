# OASIS-core

OASIS-core is a Go module that provides a simple module loader for registering services by name and an in-memory gRPC client registry. It enables dynamic service registration and inter-module communication over gRPC.

## Requirements

- Go 1.21+

## Project Structure

```text
OASIS-core/
├── core/
│   └── registry.go   # ModuleLoader, GRPCClientRegistry interface, and InMemoryGRPCClientRegistry
├── go.mod            # Module definition
└── README.md         # Documentation
```

## Installation

You can require this module in your project by adding it to your `go.mod`. Locally you might replace it like:
```go
replace OASIS-core => ../OASIS-core
```

Or just fetch it using:

```bash
go get OASIS-core
```

## Usage

### 1. Registering a Module

Define a module by implementing the `Module` interface:

```go
package main

import (
	"fmt"
	"OASIS-core/core"
)

type MyService struct{}

func (s *MyService) Name() string {
	return "my-service"
}

func (s *MyService) Start() error {
	fmt.Println("Starting my-service...")
	return nil
}

func (s *MyService) Stop() error {
	fmt.Println("Stopping my-service...")
	return nil
}

func main() {
	loader := core.NewModuleLoader()
	
	myService := &MyService{}
	
	// Register the module
	err := loader.Register(myService)
	if err != nil {
		panic(err)
	}
	
	// Retrieve and start the module
	m, err := loader.Get("my-service")
	if err == nil {
		m.Start()
	}
}
```

### 2. Calling Another Module over gRPC

Use the `GRPCClientRegistry` to register and retrieve gRPC connections to other modules:

```go
package main

import (
	"log"
	"OASIS-core/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	registry := core.NewInMemoryGRPCClientRegistry()
	
	// Connect to another module's gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	
	// Register the client connection for "target-module"
	err = registry.RegisterClient("target-module", conn)
	if err != nil {
		log.Fatalf("Failed to register client: %v", err)
	}
	
	// Retrieve the connection when needed
	clientConn, err := registry.GetClient("target-module")
	if err != nil {
		log.Fatalf("Failed to retrieve client: %v", err)
	}
	
	// Use clientConn to call standard gRPC methods, e.g.:
	// targetClient := pb.NewTargetServiceClient(clientConn)
	// resp, err := targetClient.DoSomething(context.Background(), &pb.Request{})
}
```
