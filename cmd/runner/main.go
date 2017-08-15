package main

import (
	"context"
	"log"
	"time"

	"github.com/caelifer/runner/component"
	"github.com/caelifer/runner/component/impl/job"
	"github.com/caelifer/runner/component/impl/task"
	"github.com/caelifer/runner/service/impl/mysql"
	"github.com/caelifer/runner/service/store"
)

func main() {
	// Create store.Service
	var storeService store.Service = mysql.New()

	// Create job component with tasks
	var j = job.New(storeService, []component.Task{
		task.New("converter-lres", "convert-stream", "-r 420x280"),
		task.New("converter-mres", "convert-stream", "-r 1280x720"),
		task.New("converter-hres", "convert-stream", "-r 1920x1080"),
	})

	// Create job's conext
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	// Execute job
	err := j.Run(ctx)
	if err != nil {
		log.Fatalf("runner: %v", err)
	}
}
