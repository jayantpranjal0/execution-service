# Execution Service

The Execution Service is a distributed job scheduling system designed to manage and execute tasks across multiple worker nodes. This service is responsible for coordinating job execution, managing worker nodes, and ensuring high availability and fault tolerance.

## Project Structure

```
execution-service
├── cmd
│   └── main.go                # Entry point of the execution service
├── internal
│   ├── coordinator
│   │   ├── coordinator.go      # Manages worker nodes and job assignments
│   │   └── worker_manager.go    # Handles the lifecycle of worker nodes
│   ├── worker
│   │   ├── worker.go           # Represents a worker node
│   │   └── job_executor.go      # Executes jobs assigned to the worker
│   ├── models
│   │   └── job.go              # Defines the Job struct
│   ├── queue
│   │   ├── queue.go            # Job queuing and dequeuing logic
│   │   └── kafka_client.go      # Interacts with Kafka for job messages
│   └── storage
│       ├── db.go               # Database connection and setup
│       └── job_repository.go    # Interacts with job storage
├── config
│   └── config.yaml             # Configuration settings for the service
├── go.mod                      # Go module definition
├── go.sum                      # Module dependency checksums
└── README.md                   # Documentation for the project
```

## Features

- **Job Submission**: Users can submit one-time or periodic jobs for execution.
- **Job Monitoring**: The system provides real-time monitoring of job status (queued, running, completed, failed).
- **Fault Tolerance**: The service is designed to handle worker node failures and reschedule jobs accordingly.
- **Scalability**: Capable of handling millions of jobs efficiently across multiple worker nodes.

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Kafka (for job queuing)
- A compatible database (e.g., PostgreSQL, MongoDB)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/execution-service.git
   cd execution-service
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Configure the service:
   Edit the `config/config.yaml` file to set up your database and Kafka connection details.

### Running the Service

To start the execution service, run the following command:
```
go run cmd/main.go
```

### API Endpoints

- **Submit Job**: `POST /jobs`
- **Get Job Status**: `GET /jobs/{job_id}`
- **Cancel Job**: `DELETE /jobs/{job_id}`
- **List Pending Jobs**: `GET /jobs?status=pending&user_id={user_id}`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.