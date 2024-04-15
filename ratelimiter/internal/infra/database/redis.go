package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr, pass string, db int) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: pass,
	})

	return &RedisRepository{
		client: client,
	}
}

func (rr *RedisRepository) GetBlock(ctx context.Context, key string) (*time.Time, error) {
	redisKey := fmt.Sprintf("b-%s", key)

	block, err := rr.client.Get(ctx, redisKey).Result()

	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting value: %w", err)
	}

	blockTime, err := time.Parse(time.RFC3339Nano, block)
	if err != nil {
		return nil, fmt.Errorf("error parsing data: %w", err)
	}

	return &blockTime, nil
}

func (rr *RedisRepository) IncrementAccessCount(ctx context.Context, key string, maxAccess int64) (bool, error) {
	err := rr.removeOldEntries(ctx, key)
	if err != nil {
		return false, err
	}

	count, err := rr.getCurrentCount(ctx, key)
	if err != nil {
		return false, err
	}

	if count >= maxAccess {
		return false, nil
	}

	err = rr.addNewEntry(ctx, key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (rr *RedisRepository) removeOldEntries(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("c-%s", key)
	now := time.Now()
	clearBefore := now.Add(-time.Second)

	pipe := rr.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(clearBefore.UnixMicro(), 10))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error removing old entries: %w", err)
	}

	return nil
}

func (rr *RedisRepository) getCurrentCount(ctx context.Context, key string) (int64, error) {
	redisKey := fmt.Sprintf("c-%s", key)

	count := rr.client.ZCard(ctx, redisKey).Val()

	return count, nil
}

func (rr *RedisRepository) addNewEntry(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("c-%s", key)
	now := time.Now()

	pipe := rr.client.Pipeline()
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now.UnixMicro()),
		Member: now.Format(time.RFC3339Nano),
	})
	pipe.Expire(ctx, redisKey, time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding new entry: %w", err)
	}
	return nil
}

func (rr *RedisRepository) AddBlock(ctx context.Context, key string, blockInMilliseconds int64) (*time.Time, error) {
	redisKey := fmt.Sprintf("b-%s", key)

	expiration := time.Duration(int64(time.Millisecond) * blockInMilliseconds)
	blockedUntil := time.Now().Add(expiration)

	err := rr.client.Set(ctx, redisKey, blockedUntil.Format(time.RFC3339Nano), expiration).Err()
	if err != nil {
		return nil, fmt.Errorf("error setting block: %w", err)
	}

	return &blockedUntil, nil
}
