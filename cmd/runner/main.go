package main

import (
	"context"
	"log"

	"github.com/caelifer/runner/component"
	store "github.com/caelifer/runner/service/store/mysqlstore"
)

func main() {
	// Create store.Service
	var storeService = store.New()

	// Create job component with tasks
	var job = component.NewJob(storeService, []component.Task{
		component.NewTask("converter-lres", "convert-stream", "-r 420x280"),
		component.NewTask("converter-mres", "convert-stream", "-r 1280x720"),
		component.NewTask("converter-hres", "convert-stream", "-r 1920x1080"),
	})

	// Execute job
	err := job.Run(context.Background())
	if err != nil {
		log.Fatalf("runner: %v", err)
	}
}
