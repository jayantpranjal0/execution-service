package node

import (
)

// NodeInterface is the interface that is implemented by both Coordinator and Worker.
// It provides common functionality that can be used in main.go.
type NodeInterface interface {
	// Start initializes and starts the node.
	Start() error

	// Stop gracefully stops the node.
	Stop() error

	// GetID returns the unique identifier of the node.
	GetID() string
}

