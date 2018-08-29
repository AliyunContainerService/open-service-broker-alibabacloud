package server

import (
	"fmt"
	"testing"
)

type P struct {
	a int
}

func (p P) Deal() error {
	if p.a == 1 {
		return nil
	} else {
		return fmt.Errorf("failed to deal with payload.")
	}
}

func TestQueueJob(t *testing.T) {
	NewJobQueue()
	if jobQueue == nil && len(jobQueue) != 0 {
		t.Fatalf("NewQueueJob create new job queue failed.")
	}

	jobChan := GetJobQueue()
	if jobChan != jobQueue {
		t.Fatalf("GetJobQueue get job queue failed.")
	}

	EnqueueJob(Job{Payload: P{a: 1}})
	if len(jobQueue) != 1 {
		t.Fatalf("QueueJob failed.")
	}
	jobQueue = nil
}
