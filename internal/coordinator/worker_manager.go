package coordinator

import (
	"sync"
	"time"
)

// Worker represents a worker node in the system.
type Worker struct {
	ID         string
	IPAddress  string
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
func (wm *WorkerManager) AddWorker(id, ipAddress string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.workers[id] = &Worker{
		ID:         id,
		IPAddress:  ipAddress,
		Status:     "active",
		LastHeartbeat: time.Now(),
	}
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
	return w.Status == "active"
}

func (w *Worker) AssignJob(job Job) {
	// Logic to assign a job to the worker
	w.Status = "busy"
	// Simulate job assignment
	time.Sleep(1 * time.Second)
	w.Status = "active"
}

func (w *Worker) IsFree() bool {
	return w.Status == "active"
}
