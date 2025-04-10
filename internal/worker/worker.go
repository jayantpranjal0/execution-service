package worker

import (
	"encoding/json"
	"execution-service/internal/database"
	"execution-service/internal/models"
	"execution-service/internal/queries"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Worker struct {
	ID      string
	Address string
	JobID  	string
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
	http.HandleFunc("/health", w.handleHealthRequest)
	http.HandleFunc("/job", w.handleJobRequest)

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
	log.Printf("Worker %s: Received job execution request", w.ID)
	if req.Method != http.MethodPost {
		http.Error(wr, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the job payload
	var jobPayload map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&jobPayload); err != nil {
		http.Error(wr, "Failed to parse job payload", http.StatusBadRequest)
		w.JobID = ""
		return
	}
	jobID, ok := jobPayload["JobID"].(string)
	if !ok {
		http.Error(wr, "Invalid job payload", http.StatusBadRequest)
		w.JobID = ""
		return
	}

	log.Printf("Worker %s: Received job: %v", w.ID, jobPayload)

	// Execute the job
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("Job execution started successfully"))
	if err := w.ExecuteJob(jobPayload); err != nil {
		// http.Error(wr, "Failed to execute job", http.StatusInternalServerError)
		if !ok {
			log.Printf("Worker %s: Invalid job_id in payload", w.ID)
			http.Error(wr, "Invalid job_id in payload", http.StatusBadRequest)
			w.JobID = ""
		}
		markJobCompleted(jobID, "error", err.Error())
		return
	}
	log.Printf("Worker %s: Job %s executed successfully with proof", w.ID, jobID)

	markJobCompleted(jobID, "success", "")
	log.Printf("Worker %s: Job %s executed successfully", w.ID, jobID)
	w.JobID = ""
}

func (w *Worker) ExecuteJob(jobPayload map[string]interface{}) error {
	// Simulate the job execution
	log.Printf("Worker %s: Executing job with payload: %v", w.ID, jobPayload)

	// Sleep for a random amount of milliseconds to simulate processing
	// Fetch the Dockerfile from the Firebase S3 bucket
	dockerFileURL := jobPayload["DockerfileReference"].(string)
	log.Print("Worker %s: Fetching Dockerfile from URL: %s", w.ID, dockerFileURL)
	// Fetch the Dockerfile from the provided URL
	resp, err := http.Get(dockerFileURL)
	if err != nil {
		log.Printf("Worker %s: Failed to fetch Dockerfile from URL: %v", w.ID, err)
		return err
	}
	defer resp.Body.Close()
	log.Printf("Worker %s: Received response from Dockerfile URL: %d", w.ID, resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Worker %s: Received non-OK response while fetching Dockerfile: %d", w.ID, resp.StatusCode)
		return err
	}
	log.Print("Worker %s: Successfully fetched Dockerfile from URL", w.ID)
	log.Print("creating temporary file for Dockerfile")
	// Save the Dockerfile to a temporary location
	tempFile, err := os.CreateTemp("", "dockerfile-*.Dockerfile")
	if err != nil {
		log.Printf("Worker %s: Failed to create temporary file for Dockerfile: %v", w.ID, err)
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Printf("Worker %s: Failed to save Dockerfile to temporary file: %v", w.ID, err)
		return err
	}

	// Close the file to ensure it's written to disk
	if err := tempFile.Close(); err != nil {
		log.Printf("Worker %s: Failed to close temporary Dockerfile: %v", w.ID, err)
		return err
	}
	log.Printf("Worker %s: Dockerfile saved to temporary file: %s", w.ID, tempFile.Name())
	// Execute the Dockerfile
	// Use the Docker CLI to build and run the Dockerfile
	dockerImageName := "job-image-" + jobPayload["JobID"].(string)

	// Build the Docker image
	buildCmd := exec.Command("docker", "build", "-t", dockerImageName, "-f", tempFile.Name(), ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	// log.Printf("Worker %s: Building Docker image %s", w.ID, dockerImageName)
	if err := buildCmd.Run(); err != nil {
		log.Printf("Worker %s: Failed to build Docker image: %v", w.ID, err)
		return err
	}

	// Run the Docker container
	runCmd := exec.Command("docker", "run", "--rm", dockerImageName)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	log.Printf("Worker %s: Running Docker container for image %s", w.ID, dockerImageName)
	if err := runCmd.Run(); err != nil {
		log.Printf("Worker %s: Failed to run Docker container: %v", w.ID, err)
		return err
	}

	log.Printf("Worker %s: Successfully executed Dockerfile", w.ID)

	log.Printf("Worker %s: Successfully fetched and saved Dockerfile to %s", w.ID, tempFile.Name())

	// Example: You could use a library like "github.com/docker/docker/client" to interact with Docker
	log.Printf("Worker %s: Job execution completed", w.ID)

	return nil
}

func markJobCompleted(job_id string, status string, error string) {
	// This function should update the job status in the database
	// You can use the database queries package to perform this operation
	log.Printf("status: %s", status)
	collection := database.GetCollection("hackathon", "executed_jobs")
	err := queries.AddEntry(collection, models.ExecutedJob{
		ID:                      primitive.NewObjectID(),
		JobID:                   job_id,
		ScheduledTime:           time.Now(),
		DockerfileReference:     "",
		ExecutionCompletionTime: time.Now(),
		Status:                  status,
		ErrorMessage:            error,
	})
	
	log.Printf("Worker: Marking job %s as completed with status: %s", job_id, status)
	if err != nil {
		log.Printf("Error updating job status: %v", err)
	}

}
