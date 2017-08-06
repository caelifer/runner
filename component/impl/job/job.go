package job

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/caelifer/runner/service/store"
	"github.com/caelifer/runner/util/generator"
	"github.com/caelifer/runner/component"
)

var timeout = time.Duration(5 * time.Second)

type job struct {
	id      string
	tasks   []component.Task
	success bool
	text    string
	store   store.Service
}

func New(store store.Service, tasks []component.Task) *job {
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
			"success", j.success,
			"err", err,
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
			res <- result{task.Name(), err}
		}()
	}
	go func() {
		wg.Wait()
		close(res)
	}()

	// Gather execution status for all executed tasks
	txt := []string{}
	for r := range res {
		if r.err != nil {
			j.success = false
			txt = append(txt, fmt.Sprintf("task '%v' failed: %v",
				r.tsk,
				r.err,
			))
		}
	}

	// Update persistent state
	j.text = strings.Join(txt, ", ")
	j.store.Update(j.id, j)

	if !j.success {
		err = fmt.Errorf("job %v failed: %v", j.id, j.text)
	}

	return
}

type result struct {
	tsk string
	err error
}
