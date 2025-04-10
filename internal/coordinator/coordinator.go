package coordinator

import (
	"context"
	"encoding/json"
	"execution-service/internal/queue"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Job struct {
	ID        string
	JobID     string
	WorkerID  string
	DockerfileReference string
	JobStatus string
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
	kafkaClient   *queue.KafkaClient
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

	// Start fetching jobs from Kafka
	go c.fetchJobsFromKafka()

	return nil
}

func (c *Coordinator) fetchJobsFromKafka() {
	for {
		// time.Sleep(10 * time.Second)
		jobMessage, err := c.kafkaClient.ConsumeMessage(context.TODO())
		if err != nil {
			c.logger.Error("Failed to fetch job from Kafka", zap.Error(err))
			continue
		}

		// Parse the job message
		var jobJson any
		if err := json.Unmarshal([]byte(jobMessage), &jobJson); err != nil {
			log.Print("Failed to unmarshal job message", err)
			continue
		}

		jobMap, ok := jobJson.(map[string]interface{})
		if !ok {
			log.Print("Failed to cast job message to map[string]interface{}")
			continue
		}
		
		job := Job{
			ID: "1",
			JobID: jobMap["job_id"].(string),
			WorkerID: "",
			DockerfileReference: jobMap["dockerfile_reference"].(string),
			JobStatus: "pending",
		}

		// // Enqueue the job into the jobQueue
		c.jobQueue <- job
		log.Print("Job enqueued", job.JobID, job.DockerfileReference)
	}
}

func NewCoordinator(config *viper.Viper) *Coordinator {
	kafkaClient := queue.NewKafkaClient(
		config.GetStringSlice("kafka.brokers"),
		config.GetString("kafka.topic"),
	)

	return &Coordinator{
		workers: InitializeWorkersFromConfig(config),
		mu:      sync.Mutex{},
		healthCheck: func() time.Duration {
			duration, err := time.ParseDuration(config.GetString("workers.heartbeat_interval"))
			if err != nil {
				panic(fmt.Sprintf("invalid duration for workers.heartbeat_interval: %v", err))
			}
			return duration
		}(),
		jobQueue:    make(chan Job, config.GetInt("workers.max_concurrent_jobs")),
		kafkaClient: kafkaClient,
	}
}

func (c *Coordinator) monitorWorkers() {
	for {
		time.Sleep(c.healthCheck)
		c.mu.Lock()
		for id, worker := range c.workers.workers {
			if !worker.IsHealthy() {
				log.Printf("Worker %s is unhealthy, removing from the list\n", id)
				c.workers.RemoveWorker(id)
			} else {
				// fmt.Printf("Worker %s is healthy\n", id)
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

func InitializeWorkersFromConfig(config *viper.Viper) WorkerManager {
	workerManager := NewWorkerManager()

	workers := config.Get("workers.list").([]interface{})

	for _, worker := range workers {
		workerMap := worker.(map[string]interface{}) // Convert to map[string]interface{}
		id := workerMap["id"].(string)
		name := workerMap["name"].(string)
		address := workerMap["address"].(string)
		// // Create a new worker and add it to the WorkerManager
		log.Print("Adding worker to manager", id, name, address)
		newWorker := Worker{
			ID:          id,
			Name:        name,
			Address:     address,
			AssignedJob: nil,
		}
		// newWorker.updateHealth()
		// newWorker.UpdateJobStatus()
		workerManager.AddWorker(&newWorker)
	}
	log.Print("Initializing workers from config")
	return *workerManager
}
