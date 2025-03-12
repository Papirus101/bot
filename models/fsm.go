package models

import (
	"context"
)

type FSM interface {
	SetState(ctx context.Context, key string, ttl uint) error
	GetState(ctx context.Context, key string) (any, error)
	DeleteState(ctx context.Context, key string) error
}
