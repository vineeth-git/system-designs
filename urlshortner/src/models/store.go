package models

import (
	"fmt"
	"errors"
	"os"
	"context"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"


)
type Store struct {
	data map[string]string
	rdb *redis.Client
}
func InitRedis() *redis.Client {
	add := os.Getenv("REDIS_ADDR")
	if add == "" {
		add = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{Addr: add})
}

func NewStore() *Store {
	
	return &Store{data: map[string]string{}, rdb: InitRedis()}
}

func (s *Store) Generate(ctx context.Context, url string) (string, error) {
	code := fmt.Sprint(uuid.New())
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()
	ok, err := s.rdb.SetNX(ctx, fmt.Sprintf("url:%s", code), url, 0).Result()
	// if _, ok := s.data[code]; ok {
	// 	return "", errors.New("Collision")
	// }
	// s.data[code] = url
	if err != nil {
		return "", err
	}
	if ok {
		return code, nil
	}
	return "", errors.New("Code Unavailable")
}

func (s *Store) Get(ctx context.Context, code string) (string, error) {
	// if fullUrl, ok := s.data[code]; ok {
	// 	return fullUrl, nil
	// }
	ctx, cancel := context.WithTimeout(ctx, 2 *time.Second)
	defer cancel()
	fullUrl, err := s.rdb.Get(ctx, fmt.Sprintf("url:%s", code)).Result()
	if err == redis.Nil {
		return "", errors.New("Not Found")
	}
	if err != nil {
		return "", errors.New("Storage Unavailable")
	}
	return fullUrl, nil
}
