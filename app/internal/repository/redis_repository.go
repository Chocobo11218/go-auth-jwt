package repository

// import (
// 	"context"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// type RedisRepository interface {
// 	Set(ctx context.Context, key string, value string, ttl time.Duration) error
// }

// type redisRepository struct {
// 	client *redis.Client
// }

// func NewRedisRepository(client *redis.Client) redisRepository {
// 	return &redisRepository{client: client}
// }

// func (r *redisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
// 	panic("unimplemented")
// }