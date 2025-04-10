package main

import (
	"execution-service/internal/coordinator"
	"execution-service/internal/node"
	"execution-service/internal/worker"
	"os"
	"os/signal"
	"syscall"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"execution-service/internal/database"
)

var logger *zap.Logger

func main() {
	// Load configuration
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")

	if err:= viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Initialize logger
	logger, _ = zap.NewProduction()
	// if err != nil {
	// 	panic(err)
	// }
	defer logger.Sync() // flushes buffer, if any
	logger.Info("Logger initialized")

	// Setup Database Connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logger.Fatal("MONGO_URI environment variable is not set")
	}
	if err := database.ConnectMongoDB(mongoURI); err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer database.DisconnectMongoDB()
	logger.Info("Connected to MongoDB")

	// Start node
	// TODO: Every node starts as a worker and then only one node becomes a coordinator through some consensus algorithm. Also, let the cluster owner decide the coordinator as a config
	logger.Info("Starting node...")
	node,err:= NewNode(viper.GetViper())
	if err != nil {
		logger.Fatal("Failed to create node", zap.Error(err))
	}

	err = node.Start()
	if err != nil {
		logger.Fatal("Failed to start node", zap.Error(err))
	}
	logger.Info("Node started", zap.String("nodeID", node.GetID()))
	
	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("Termination signal received, shutting down...")

	// Implement graceful shutdown logic here
	node.Stop()
	logger.Info("Coordinator stopped")
}

// NewNode creates a new Node instance based on the configuration provided by viper.
// It initializes either a Worker or Coordinator node based on the "node.type" configuration.
func NewNode(config *viper.Viper) (node.NodeInterface, error) {
	nodeType := config.GetString("node.type")
	logger.Info("Node type", zap.String("type", nodeType))
	switch nodeType {
		case "worker":
			return worker.NewWorker(config), nil
		case "coordinator":
			return coordinator.NewCoordinator(config), nil
		default:
			return nil, fmt.Errorf("unknown node type: %s", nodeType)
	}
}