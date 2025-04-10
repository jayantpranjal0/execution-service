package worker

import "net/http"


func (w *Worker) handleHealthRequest(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("OK"))
}

func (w *Worker) handleJobRequest(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("10"))
}