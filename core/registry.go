package core

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

// Module represents a generic service module that can be registered.
// It defines a standard interface for starting and stopping a module.
type Module interface {
	// Name returns the unique name of the module.
	Name() string
	// Start initializes and starts the module.
	Start() error
	// Stop shuts down the module.
	Stop() error
}

// ModuleLoader is responsible for registering and managing modules by name.
// It provides thread-safe access to a shared collection of modules.
type ModuleLoader struct {
	mu      sync.RWMutex
	modules map[string]Module
}

// NewModuleLoader creates a new initialized ModuleLoader.
func NewModuleLoader() *ModuleLoader {
	return &ModuleLoader{
		modules: make(map[string]Module),
	}
}

// Register registers a new module by its name.
// Returns an error if a module with the same name is already registered.
func (l *ModuleLoader) Register(m Module) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	name := m.Name()
	if _, exists := l.modules[name]; exists {
		return fmt.Errorf("module %s is already registered", name)
	}
	l.modules[name] = m
	return nil
}

// Get retrieves a registered module by its name.
// Returns an error if the module is not found.
func (l *ModuleLoader) Get(name string) (Module, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	m, exists := l.modules[name]
	if !exists {
		return nil, fmt.Errorf("module %s not found", name)
	}
	return m, nil
}

// GRPCClientRegistry defines an interface for components that need to
// discover and retrieve gRPC clients to communicate with other modules.
type GRPCClientRegistry interface {
	// RegisterClient registers a new gRPC client connection for a specific module name.
	RegisterClient(moduleName string, conn *grpc.ClientConn) error
	
	// GetClient retrieves a gRPC client connection for a specific module name.
	GetClient(moduleName string) (*grpc.ClientConn, error)
}

// InMemoryGRPCClientRegistry is a minimal in-memory map implementation of GRPCClientRegistry.
// This is useful for storing active gRPC connections to other services.
type InMemoryGRPCClientRegistry struct {
	mu      sync.RWMutex
	clients map[string]*grpc.ClientConn
}

// NewInMemoryGRPCClientRegistry creates a new initialized InMemoryGRPCClientRegistry.
func NewInMemoryGRPCClientRegistry() *InMemoryGRPCClientRegistry {
	return &InMemoryGRPCClientRegistry{
		clients: make(map[string]*grpc.ClientConn),
	}
}

// RegisterClient stores the module name and its gRPC client connection.
// Returns an error if a client for the given module name is already registered.
func (r *InMemoryGRPCClientRegistry) RegisterClient(moduleName string, conn *grpc.ClientConn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[moduleName]; exists {
		return fmt.Errorf("gRPC client for module %s is already registered", moduleName)
	}
	r.clients[moduleName] = conn
	return nil
}

// GetClient retrieves the gRPC client connection for the specified module name.
// Returns an error if no client is found for the given module name.
func (r *InMemoryGRPCClientRegistry) GetClient(moduleName string) (*grpc.ClientConn, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conn, exists := r.clients[moduleName]
	if !exists {
		return nil, fmt.Errorf("gRPC client for module %s not found", moduleName)
	}
	return conn, nil
}
