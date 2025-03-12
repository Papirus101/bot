package fsm

import (
	"context"
	"time"

	"github.com/go-telegram/bot/models"
	"github.com/jftuga/TtlMap"
	"github.com/redis/go-redis/v9"
)

type MemoryFSM struct {
	data *TtlMap.TtlMap[string]
}

func NewMemoryFSM() models.FSM {
	return &MemoryFSM{
		data: TtlMap.New[string](time.Duration(1), 1, time.Duration(1), false),
	}
}

func (mfsm *MemoryFSM) SetState(ctx context.Context, key string, ttl uint) error {
	mfsm.data.Put(key, 1)
	return nil
}

func (mfsm *MemoryFSM) GetState(ctx context.Context, key string) (any, error) {
	v := mfsm.data.Get(key)
	return v, nil
}

func (mfsm *MemoryFSM) DeleteState(ctx context.Context, key string) error {
	mfsm.data.Delete(key)
	return nil
}

type RedisFSM struct {
	client *redis.Client
}

func NewRedisFSM(client *redis.Client) models.FSM {
	return &RedisFSM{
		client: client,
	}
}

func (rfsm *RedisFSM) SetState(ctx context.Context, key string, ttl uint) error {
	return rfsm.client.Set(ctx, key, 1, time.Duration(ttl)*time.Second).Err()
}

func (rfsm *RedisFSM) GetState(ctx context.Context, key string) (any, error) {
	val, err := rfsm.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (rfsm *RedisFSM) DeleteState(ctx context.Context, key string) error {
	return rfsm.client.Del(ctx, key).Err()
}
