// Package database deals with db connection
package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, address string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
