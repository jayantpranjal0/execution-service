package queue

import (
	"errors"
	"sync"
)

// Job represents a job in the queue.
type Job struct {
	ID       string
	Payload  interface{}
}

// Queue defines the interface for a job queue.
type Queue interface {
	Enqueue(job Job) error
	Dequeue() (Job, error)
	IsEmpty() bool
}

// InMemoryQueue is an in-memory implementation of the Queue interface.
type InMemoryQueue struct {
	jobs []Job
	mu   sync.Mutex
}

// NewInMemoryQueue creates a new instance of InMemoryQueue.
func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		jobs: make([]Job, 0),
	}
}

// Enqueue adds a job to the queue.
func (q *InMemoryQueue) Enqueue(job Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.jobs = append(q.jobs, job)
	return nil
}

// Dequeue removes and returns the next job from the queue.
func (q *InMemoryQueue) Dequeue() (Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return Job{}, errors.New("queue is empty")
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job, nil
}

// IsEmpty checks if the queue is empty.
func (q *InMemoryQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.jobs) == 0
}