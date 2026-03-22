package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"

	"OASIS-core/core"
	"oasis-data/adapters/sqlite"
	"zhangbolong/OASIS-hr/hr"
)

func main() {
	// 1. Initialize Core Loader and Registry
	loader := core.NewModuleLoader()
	clientRegistry := core.NewInMemoryGRPCClientRegistry()

	fmt.Println("Initializing Central OASIS Server...")

	// 2. Initialize Data Adapters
	employeeAdapter := sqlite.NewSqliteEmployeeAdapter()
	deptAdapter := sqlite.NewSqliteDepartmentAdapter()

	// 3. Instantiate Modules
	hrModule := hr.NewHRModule(employeeAdapter, deptAdapter)

	// 4. Register and Start Modules
	modules := []core.Module{hrModule}

	for _, mod := range modules {
		if err := loader.Register(mod); err != nil {
			fmt.Printf("Failed to register module %s: %v\n", mod.Name(), err)
			os.Exit(1)
		}

		if err := mod.Start(); err != nil {
			fmt.Printf("Failed to start module %s: %v\n", mod.Name(), err)
			os.Exit(1)
		}
	}

	// 5. Setup internal gRPC client connections for inter-module networking
	// Right now we only have OASIS-hr running on :50051
	if conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		_ = clientRegistry.RegisterClient("OASIS-hr", conn)
		fmt.Println("Registered OASIS-hr gRPC client within core registry.")
	}

	fmt.Println("Central OASIS Server running. Press Ctrl+C to stop.")

	// 6. Wait for Termination Signal to smoothly shut down
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	fmt.Println("\nShutting down Central OASIS Server...")
	for _, mod := range modules {
		mod.Stop()
	}
}
