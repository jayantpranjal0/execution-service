package worker

import (
	"math/rand"
	"log"
	"time"
)

type Config struct {
}

type Worker struct {
	ID          string
}

func NewWorker() *Worker {
	return &Worker{
		// ID:          id,
		// JobExecutor: NewJobExecutor(),
	}
}

func (w Worker) Start() error {
	log.Printf("Worker %s: Starting", w.ID)
	go func() {
		for {
			err := w.ExecuteJob()
			if err != nil {
				log.Printf("Worker %s: Error executing job: %v", w.ID, err)
			}
			time.Sleep(1 * time.Second) // Add a delay between job executions
		}
	}()
	return nil
}

func (w Worker) Stop() error {
	log.Printf("Worker %s: Stopping", w.ID)
	return nil
}

func (w Worker) GetID() string {
	return w.ID
}

func (w *Worker) ExecuteJob() error {
	// Simulate the job execution
	log.Printf("Worker %s: Executing job", w.ID)
	// Sleep for a random amount of milliseconds
	randomDuration := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(randomDuration)
	// TODO: Write code to run dockerfile
	return nil
}
