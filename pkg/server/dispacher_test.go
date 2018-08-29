package server

import (
	"testing"
)

func TestNewDispatcher(t *testing.T) {
	NewJobQueue()

	dispatcher := NewDispatcher()
	if dispatcher == nil || dispatcher.WorkerPool == nil {
		t.Fatalf("NewDispatcher create dispatcher failed.")
	}

	dispatcher = nil
}
