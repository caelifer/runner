package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caelifer/runner/component"
	"github.com/caelifer/runner/service/store"
	"github.com/caelifer/runner/util/generator"
)

var timeout = time.Duration(5 * time.Second)

type job struct {
	id      string
	tasks   []component.Task
	success bool
	text    string
	store   store.Service
	logger  *log.Logger
}

type logrec struct {
	Component string `json:"component"`
	ID        string `json:"id"`
	Status    string `json:"status"`
	Success   bool   `json:"success,omitempty"`
	Error     string `json:"error,omitempty"`
	Duration  string `json:"duration,omitempty"`
}

func (l logrec) String() string {
	out, _ := json.Marshal(&l)
	return string(out)
}

func New(store store.Service, tasks ...component.Task) *job {
	j := &job{
		id:     generator.NewID(),
		tasks:  tasks,
		store:  store,
		logger: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile),
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
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		j.logger.Printf("%v",
			logrec{
				Component: "job",
				ID:        j.id,
				Status:    "finished",
				Success:   j.success,
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	j.logger.Printf("%v",
		logrec{
			Component: "job",
			ID:        j.id,
			Status:    "started",
		},
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
