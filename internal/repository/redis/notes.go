package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/countenum404/Veksel/internal/types"
	"github.com/countenum404/Veksel/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(addr string, password string, db int) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	return &RedisRepository{client: rdb, ctx: ctx}
}

func (rr *RedisRepository) GetNotesByUserId(userId int64) ([]types.Note, error) {
	res, err := rr.client.Get(rr.ctx, fmt.Sprintf("%v", userId)).Result()
	if err != nil {
		return nil, err
	}

	var notes []types.Note
	json.Unmarshal([]byte(res), &notes)

	return notes, nil
}

func (rr *RedisRepository) PutNotes(userId int64, notes []types.Note) error {
	val, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	err = rr.client.Set(rr.ctx, fmt.Sprintf("%v", userId), val, time.Minute).Err()
	if err != nil {
		logger.GetLogger().Err(err)
		return err
	}
	return nil
}
