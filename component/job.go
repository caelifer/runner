package component

import "context"

type Job interface {
	Run(ctx context.Context) error
	Success() bool
}
