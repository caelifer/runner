package component

import "context"

type Task interface {
	Execute(ctx context.Context) error
	Name() string
	ID() string
	Success() bool
}
