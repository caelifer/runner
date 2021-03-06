package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/caelifer/runner/service/generator"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type task struct {
	id     string
	name   string
	cmd    string
	args   []string
	err    error
	logger *log.Logger
}

type logrec struct {
	Component string `json:"component"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Operation string `json:"operation"`
	Cmd       string `json:"cmd"`
	Error     string `json:"error,omitempty"`
	Duration  string `json:"duration"`
}

func (l logrec) String() string {
	out, _ := json.Marshal(&l)
	return string(out)
}

func New(name string, cmd string, args ...string) *task {
	return &task{
		id:     generator.NewID(),
		name:   name,
		cmd:    cmd,
		args:   args,
		logger: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile),
	}
}

func (t *task) Execute(ctx context.Context) (err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		t.logger.Printf("%v",
			logrec{
				Component: "task",
				ID:        t.id,
				Name:      t.name,
				Operation: "execute",
				Cmd:       strings.Join(append([]string{t.cmd}, t.args...), " "),
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	// cmd := exec.CommandContext(ctx, t.cmd, t.args...)
	pause := fmt.Sprintf("%.2f", (time.Duration(500+rand.Intn(5000)) * time.Millisecond).Seconds())
	// Simulate work
	cmd := exec.CommandContext(ctx, "sleep", pause)
	// Use our own STDERR for task's diagnostic messages
	cmd.Stderr = os.Stderr
	// Run external command
	if err = cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			err = errors.New("execution timed out")
		}
	}

	// Update state
	t.err = err

	return
}

func (t *task) Name() string {
	return t.name
}

func (t *task) String() string {
	return fmt.Sprintf("%v: %q", t.name,
		strings.Join(append([]string{t.cmd}, t.args...), " "))
}

func (t *task) ID() string {
	return t.id
}

func (t *task) Success() bool {
	return t.err == nil
}
