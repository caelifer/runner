package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/caelifer/runner/component/job"
	"github.com/caelifer/runner/component/task"
	"github.com/caelifer/runner/service/store/memory"
)

func main() {
	logger := log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	// Create store.Service
	var storeService = memory.New()
	// Create job component with tasks
	var j = job.New(
		storeService,
		task.New("low-res", "convert-stream", "-r 420x280"),
		task.New("mid-res", "convert-stream", "-r 1280x720"),
		task.New("hi-res", "convert-stream", "-r 1920x1080"),
	)
	// Create job's context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	// Execute job
	err := j.Run(ctx)
	if err != nil {
		logger.Fatalf("runner: %v", err)
	}
}
