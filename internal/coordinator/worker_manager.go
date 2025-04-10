package coordinator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"log"	
)

// Worker represents a worker node in the system.
type Worker struct {
	ID         string
	Name       string
	Address    string
	Status     string
	AssignedJob *Job
	LastHeartbeat time.Time
}

// WorkerManager manages the lifecycle of worker nodes.
type WorkerManager struct {
	workers map[string]*Worker
	mu      sync.Mutex
}

// NewWorkerManager creates a new WorkerManager instance.
func NewWorkerManager() *WorkerManager {
	return &WorkerManager{
		workers: make(map[string]*Worker),
	}
}

// AddWorker adds a new worker to the manager.
func (wm *WorkerManager) AddWorker(w *Worker) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.workers[w.ID] = w
}

// RemoveWorker removes a worker from the manager.
func (wm *WorkerManager) RemoveWorker(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(wm.workers, id)
}

// UpdateWorkerStatus updates the status of a worker.
func (wm *WorkerManager) UpdateWorkerStatus(id, status string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if worker, exists := wm.workers[id]; exists {
		worker.Status = status
		worker.LastHeartbeat = time.Now()
	}
}

// GetActiveWorkers returns a list of active workers.
func (wm *WorkerManager) GetActiveWorkers() []*Worker {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	activeWorkers := make([]*Worker, 0)
	for _, worker := range wm.workers {
		if worker.Status == "active" {
			activeWorkers = append(activeWorkers, worker)
		}
	}
	return activeWorkers
}

// CheckWorkerHealth checks the health of workers and removes inactive ones.
func (wm *WorkerManager) CheckWorkerHealth(timeout time.Duration) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	for _, worker := range wm.workers {
		if time.Since(worker.LastHeartbeat) > timeout {
			worker.Status = "inactive"
		}
	}
}

func (w *Worker) IsHealthy() bool {
	resp, err := http.Get(w.Address + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (w *Worker) AssignJob(job Job) {
	// Logic to assign a job to the worker
	w.Status = "busy"
	log.Print("Assigning job to worker", w.ID, job.ID)
	reqBody, err := json.Marshal(job)
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}

	resp, err := http.Post(w.Address+"/execute", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// panic("Failed to assign job") // Handle error appropriately in production code
		log.Printf("Failed to assign job to worker %s: %s", w.ID, resp.Status, resp.Body)
		return
	} else {
		log.Printf("Job assigned to worker %s successfully", w.ID)
	}

	w.Status = "active"
}

func (w *Worker) IsFree() bool {
	resp, err := http.Get(w.Address + "/job")
	if err != nil {
		log.Printf("Error checking job status for worker %s: %v", w.ID, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var job any
		if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
			log.Printf("Error decoding job response for worker %s: %v", w.ID, err)
			return false
		}
		jobMap, ok := job.(map[string]interface{})
		if !ok {
			log.Printf("Invalid job response format for worker %s", w.ID)
			return false
		}
		jobID, ok := jobMap["JobID"].(string)
		if !ok {
			log.Printf("Invalid job ID format for worker %s", w.ID)
			return false
		}
		if jobID == "" {
			return true // No job assigned
		}
		return false // Job is assigned
	}
	log.Printf("Error checking job status for worker %s: %v", w.ID, err)
	return false
}


// func (w *Worker) updateHealth() error {
// 	resp, err := http.Get(w.Address + "/health")
// 	if err != nil {
// 		w.Status = "inactive"
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		w.Status = "active"
// 	} else {
// 		w.Status = "inactive"
// 	}
// 	return nil
// }

func (w* Worker) UpdateJobStatus() error {
	resp, err := http.Get(w.Address + "/job")
	if err != nil {

	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Logic to update the assigned job
		// For example, parse the response body to get job details
		var job Job
		if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
			return err
		}
		w.AssignedJob = &job
	} else {
		w.AssignedJob = nil
	}
	return nil
}