package server

import (
	"os"
	"strconv"

	"github.com/golang/glog"
)

var (
	WorkerMax      = os.Getenv("MAX_WORKERS")
	WorkerQueueMax = os.Getenv("MAX_WORKER_QUEUE")

	WorkerMaxDefault      = 10
	WorkerQueueMaxDefault = 10
)

// Dispatcher dispatch incoming jobs to idle worker's job channel
type Dispatcher struct {
	// A pool of each worker's job channel that are registered with the dispatcher
	WorkerPool chan chan Job
}

func NewDispatcher() *Dispatcher {
	workerMax, err := strconv.Atoi(WorkerMax)
	if err != nil {
		glog.Infof("Environment variable MAX_WORKERS is not set.\n")
		workerMax = WorkerQueueMaxDefault
	}

	pool := make(chan chan Job, workerMax)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	workerQueueMax, err := strconv.Atoi(WorkerQueueMax)
	if err != nil {
		glog.Infof("Environment variable MAX_WORKER_QUEUE is not set.\n")
		workerQueueMax = WorkerQueueMaxDefault
	}

	// starting workers
	workerMax, err := strconv.Atoi(WorkerMax)
	if err != nil {
		workerMax = WorkerMaxDefault
	}

	for i := 0; i < workerMax; i++ {
		worker := NewWorker(d.WorkerPool, workerQueueMax)
		worker.Start()
	}
	glog.Infof("%d workers are started.", workerMax)
	// start to dispatch jobs
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	glog.Info("Start to dispatch request")
	for {
		select {
		// blocking until job comes in
		case job := <-GetJobQueue():
			// a job has been received
			go func(job Job) {
				// try to obtain a worker's job channel.
				// this will block until any worker is registered in the pool
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
