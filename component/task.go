package component

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Task interface {
	Execute(ctx context.Context) error
	Name() string
}

type task struct {
	name string
	cmd  string
	args []string
	err  error
}

func NewTask(name string, cmd string, args ...string) Task {
	return &task{
		name: name,
		cmd:  cmd,
		args: args,
	}
}

func (t *task) Execute(ctx context.Context) (err error) {
	defer func(t0 time.Time) {
		logger.Log(
			"component", "task",
			"name", t.name,
			"operation", "execute",
			"cmd", strings.Join(append([]string{t.cmd}, t.args...), " "),
			"error", err,
			"duration", time.Since(t0),
		)
	}(time.Now())

	// cmd := exec.CommandContext(ctx, t.cmd, t.args...)
	pause := fmt.Sprintf("%.2f", (time.Duration(500+rand.Intn(5000)) * time.Millisecond).Seconds())
	// Simulate work
	cmd := exec.CommandContext(ctx, "sleep", pause)
	cmd.Stderr = os.Stderr
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
