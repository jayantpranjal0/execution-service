package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ScheduledJob represents the schema for a scheduled job
type ScheduledJob struct {
    ID                 primitive.ObjectID `bson:"_id,omitempty"`            // MongoDB ObjectID
    JobID              string             `bson:"job_id"`                  // Unique Job ID
    DockerfileReference string            `bson:"dockerfile_reference"`    // Reference to the Dockerfile
    ScheduledTime      time.Time          `bson:"scheduled_time"`          // Time when the job is scheduled
    CronExpression     string             `bson:"cronexpression"`          // Cron expression for recurring jobs
}

type ExecutedJob struct {
    ID                 primitive.ObjectID `bson:"_id,omitempty"`            // MongoDB ObjectID
    JobID              string             `bson:"job_id"`                  // Unique Job ID
    DockerfileReference string            `bson:"dockerfile_reference"`    // Reference to the Dockerfile
    ScheduledTime      time.Time          `bson:"scheduled_time"`          // Time when the job is scheduled
    ExecutionCompletionTime time.Time `bson:"execution_completion_time"` // Time when the job execution is completed
    Status             string             `bson:"status"`                  // Status of the job (e.g., "completed", "failed")
    ErrorMessage       string             `bson:"error_message"`           // Error message if the job failed
}