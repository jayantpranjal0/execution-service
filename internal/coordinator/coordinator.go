package coordinator

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Job struct {
	ID       string
	Task     string
	WorkerID string
}

type Config struct {
	JobQueueSize        int
	WorkerTimeout       time.Duration
	HealthCheckInterval time.Duration
}

type Coordinator struct {
	logger        *zap.Logger
	workers       WorkerManager
	mu            sync.Mutex
	healthCheck   time.Duration
	jobQueue      chan Job
	workerTimeout time.Duration
}

func (c *Coordinator) Stop() error {
	panic("unimplemented")
}

func (c *Coordinator) GetID() string {
	return "coordinator"
}

func (c *Coordinator) Start() error {
	fmt.Printf("Coordinator started\n")
	go c.monitorWorkers()
	// c.logger.Info("Coordinator started")
	// setUpWorkersFromConfig()
	return nil
}

func NewCoordinator(config *viper.Viper) *Coordinator {
	return &Coordinator{
		workers:     WorkerManager{
			workers: make(map[string]*Worker),
		},
		mu:          sync.Mutex{},
		// healthCheck: coordinatorConfig.HealthCheckInterval,
		// jobQueue:	make(chan Job, coordinatorConfig.JobQueueSize),
		// workerTimeout: coordinatorConfig.WorkerTimeout,
	}
}

func (c *Coordinator) monitorWorkers() {
	for {
		time.Sleep(c.healthCheck)
		c.mu.Lock()
		for id, worker := range c.workers.workers {
			if !worker.IsHealthy() {
				fmt.Print("Worker %s is unhealthy, removing from the list\n", id)
				c.workers.RemoveWorker(id)
			} else {
				fmt.Print("Worker %s is healthy\n", id)
			}
			if worker.IsFree() {
				select {
				case job := <-c.jobQueue:
					worker.AssignJob(job)
				default:
					// No job available in the queue
				}
			}
		}
		c.mu.Unlock()
	}
}
