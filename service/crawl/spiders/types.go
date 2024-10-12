package spiders

import "context"

type Spider interface {
	Start(ctx context.Context)
}
