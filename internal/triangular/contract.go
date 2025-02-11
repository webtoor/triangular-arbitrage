package triangular

import "context"

type Resolve interface {
	Start(ctx context.Context) error
}
