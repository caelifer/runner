package component

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/caelifer/runner/util/generator"
	"github.com/caelifer/runner/service/store"
)

var timeout = time.Duration(3 * time.Second)

type Job interface {
	Run(ctx context.Context) error
	Success() bool
}

type job struct {
	id      string
	tasks   []Task
	success bool
	text    string
	store   store.Service
}

func NewJob(tasks []Task, store store.Service) Job {
	j := &job{
		id:    generator.NewID(),
		tasks: tasks,
		store: store,
	}

	j.store.Create(j)

	return j
}

func (j *job) ID() string {
	return j.id
}

func (j *job) Success() bool {
	return j.success
}

func (j *job) Run(ctx context.Context) (err error) {
	defer func(t0 time.Time) {
		logger.Log(
			"component", "job",
			"id", j.id,
			"status", "finished",
			"error", err,
			"duration", time.Since(t0),
		)
	}(time.Now())

	logger.Log(
		"component", "job",
		"id", j.id,
		"status", "started",
	)

	j.success = true // assume all is going to be well

	var res = make(chan result, len(j.tasks))
	var wg sync.WaitGroup

	wg.Add(len(j.tasks))
	for _, task := range j.tasks {
		task := task
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			err := task.Execute(ctx)
			res <- result{task.String(), err}
		}()
	}
	go func() {
		wg.Wait()
		close(res)
	}()

	// Update persistent state
	for r := range res {
		if r.err != nil {
			j.success = false
			j.text += fmt.Sprintf(" task %q failed: %v;",
				r.tsk,
				r.err,
			)
		}
	}
	j.store.Update(j.id, j)

	if !j.success {
		err = fmt.Errorf("job %v failed", j.id)
	}

	return
}

type result struct {
	tsk string
	err error
}
