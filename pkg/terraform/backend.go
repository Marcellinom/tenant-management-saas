package terraform

import (
	"context"
)

type TfBackend interface {
	// Init butuh terraform executable dalam konteks dengan key "terraform"
	Init(ctx context.Context) error
}
