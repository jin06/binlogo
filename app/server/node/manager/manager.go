package manager

import "context"

type Manager interface {
	Start(ctx context.Context) chan error
}
