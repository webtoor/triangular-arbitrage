package triangular

import "context"

type Resolverer interface {
	Start(ctx context.Context) error
}
