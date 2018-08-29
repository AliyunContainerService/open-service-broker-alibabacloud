package server

import (
	"testing"
	"time"
	//"fmt"
)

func TestStart(t *testing.T) {
	workerPool := make(chan chan Job, 10)
	worker := NewWorker(workerPool, 10)
	if len(worker.WorkerPool) != 0 || (worker.JobChannel == nil) || worker.quit == nil {
		t.Fatalf("NewWorker create new worker failed.")
	}
	worker.Start()
	time.Sleep(1 * time.Second)
	if len(worker.WorkerPool) != 1 {
		t.Fatalf("Start worker failed.")
	}

	// pick a worker, send it a correct payload to handle
	// the worker should be put back to pool after handling
	w_job_channel := <-worker.WorkerPool
	w_job_channel <- Job{Payload: P{a: 1}}
	time.Sleep(time.Second)
	if len(worker.WorkerPool) != 1 {
		t.Fatalf("Worker return nil to keep runing failed.")
	}

	// worker should deal with incorrect payload, and back to pool
	w_job_channel = <-worker.WorkerPool
	w_job_channel <- Job{Payload: P{a: 0}}
	time.Sleep(time.Second)
	if len(worker.WorkerPool) != 1 {
		t.Fatalf("Worker return error to keep runing failed.")
	}

	time.Sleep(time.Second)
}
