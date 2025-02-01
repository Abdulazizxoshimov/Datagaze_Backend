package redis

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/storage"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/interfaces"
)



type cache struct {
	rdb *redis.RedisDB
}
func NewRedis(rdb *redis.RedisDB)interfaces.Redis{
	return &cache{
		rdb: rdb,
	}
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	byteData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.rdb.Client.Set(ctx, key, string(byteData), expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.rdb.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	
	return []byte(data), nil
}

func (c *cache) Del(ctx context.Context, key string) error {
	err := c.rdb.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
