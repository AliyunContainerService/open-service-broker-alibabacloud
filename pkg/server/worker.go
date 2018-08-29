package server

import (
	"github.com/golang/glog"
)

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, maxQueue int) *Worker {
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job, maxQueue),
		quit:       make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			//register the current worker into the dispatch's worker pool
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				// Worker receives a request
				if err := job.Payload.Deal(); err != nil {
					/* TODO: need to handle failed request,
					   maybe enqueue it back to job queue with retry limits
					*/
					glog.Infof("Error to deal request payload: %v\n", err.Error())
				}
				//case <-w.quit:
				//	// Received stop signal
				//	return
			}
		}
	}()
}

//// Stop set signal for worker to stop dealing with requests
//func (w *Worker) Stop() {
//	go func() {
//		w.quit <- true
//	}()
//}
