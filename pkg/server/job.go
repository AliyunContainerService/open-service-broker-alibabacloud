package server

import (
	"os"
	"strconv"

	"github.com/golang/glog"
)

var (
	JobQueueMax = os.Getenv("MAX_JOB_QUEUE")

	JobQueueMaxDefault = 10
)

type JobPayload interface {
	Deal() error
}

type Job struct {
	Payload JobPayload
}

// jobQueue a queue to buffer incoming job
var jobQueue chan Job

func NewJobQueue() {
	jobQueueMax, err := strconv.Atoi(JobQueueMax)
	if err != nil {
		glog.Infof("Environment variable MAX_JOB_QUEUE is not set.\n")
		jobQueueMax = JobQueueMaxDefault
	}
	jobQueue = make(chan Job, jobQueueMax)
	glog.Infof("Create a job queue.\n")
	return
}

func EnqueueJob(job Job) {
	jobQueue <- job
	return
}

func GetJobQueue() chan Job {
	return jobQueue
}
