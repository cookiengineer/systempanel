package model

import "context"

type Observer func()

type Model interface {
	Refresh(ctx context.Context) error
	Observe(fn Observer) func()
}
