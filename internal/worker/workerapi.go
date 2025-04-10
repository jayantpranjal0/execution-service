package worker

import (
	"encoding/json"
	"net/http"
)


func (w *Worker) handleHealthRequest(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("OK"))
}

func (w *Worker) handleJobRequest(wr http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		jobID := w.JobID // Assuming getJobID() is a method that retrieves the current job ID
		response := map[string]string{"JobID": jobID}
		wr.Header().Set("Content-Type", "application/json")
		wr.WriteHeader(http.StatusOK)
		json.NewEncoder(wr).Encode(response)
		return
	}
	wr.WriteHeader(http.StatusMethodNotAllowed)
	wr.Write([]byte("Method not allowed"))
}