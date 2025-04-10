package worker

import (
	"encoding/json"
	"execution-service/internal/database"
	"execution-service/internal/models"
	"execution-service/internal/queries"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	Address string // Address where the worker listens (e.g., ":8080")
}

type Worker struct {
	ID      string
	Address string
}

func NewWorker(config *viper.Viper) *Worker {
	return &Worker{
		ID:      config.GetString("node.id"),
		Address: config.GetString("node.address"),
	}
}

func (w *Worker) Start() error {
	log.Printf("Worker %s: Starting on %s", w.ID, w.Address)

	// Define HTTP handlers
	http.HandleFunc("/execute", w.handleExecuteJob)

	// Start the HTTP server
	go func() {
		if err := http.ListenAndServe(w.Address, nil); err != nil {
			log.Fatalf("Worker %s: Failed to start HTTP server: %v", w.ID, err)
		}
	}()

	return nil
}

func (w *Worker) Stop() error {
	log.Printf("Worker %s: Stopping", w.ID)
	// TODO: Implement graceful shutdown logic if needed
	return nil
}

func (w *Worker) GetID() string {
	return w.ID
}

func (w *Worker) handleExecuteJob(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(wr, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the job payload
	var jobPayload map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&jobPayload); err != nil {
		http.Error(wr, "Failed to parse job payload", http.StatusBadRequest)
		return
	}
	jobID, ok := jobPayload["job_id"].(string)

	log.Printf("Worker %s: Received job: %v", w.ID, jobPayload)

	// Execute the job
	if err := w.ExecuteJob(jobPayload); err != nil {
		// http.Error(wr, "Failed to execute job", http.StatusInternalServerError)
		if !ok {
			log.Printf("Worker %s: Invalid job_id in payload", w.ID)
			http.Error(wr, "Invalid job_id in payload", http.StatusBadRequest)
			return
		}
		markJobCompleted(jobID, "error", err.Error())
		return
	}

	markJobCompleted(jobID, "success", "")
}

func (w *Worker) ExecuteJob(jobPayload map[string]interface{}) error {
	// Simulate the job execution
	log.Printf("Worker %s: Executing job with payload: %v", w.ID, jobPayload)

	// Sleep for a random amount of milliseconds to simulate processing
	randomDuration := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(randomDuration)

	// TODO: Write code to run a Dockerfile
	// Example: You could use a library like "github.com/docker/docker/client" to interact with Docker
	log.Printf("Worker %s: Job execution completed", w.ID)

	return nil
}

func markJobCompleted(job_id string, status string, error string) {
	// This function should update the job status in the database
	// You can use the database queries package to perform this operation
	collection := database.GetCollection("hackathons", "executed_jobs")
	err := queries.AddEntry(collection, models.ExecutedJob{
		ID:                      primitive.NewObjectID(),
		JobID:                   job_id,
		ScheduledTime:           time.Now(),
		DockerfileReference:     "",
		ExecutionCompletionTime: time.Now(),
		Status:                  status,
		ErrorMessage:            error,
	})


	if err != nil {
		log.Printf("Error updating job status: %v", err)
	}

}
