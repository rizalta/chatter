package database

import (
	"chatter/server/internal/user"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound              = errors.New("key not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type UserRepo struct {
	db *redis.Client
}

func NewUserRepo(db *redis.Client) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser(ctx context.Context, u *user.User) error {
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now().UTC()

	userKey := fmt.Sprintf("user:%s", u.ID)
	userNameKey := fmt.Sprintf("username:%s", u.Username)

	err := r.db.Watch(ctx, func(tx *redis.Tx) error {
		exists, err := tx.Exists(ctx, userNameKey).Result()
		if err != nil {
			return err
		}
		if exists == 1 {
			return ErrUsernameAlreadyExists
		}

		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			p.HSet(ctx, userKey, map[string]any{
				"id":         u.ID,
				"username":   u.Username,
				"password":   u.Password,
				"created_at": u.CreatedAt.Format(time.RFC3339),
			})
			p.Set(ctx, userNameKey, u.ID, 0)
			return nil
		})

		return err
	}, userNameKey)

	return err
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	usernameKey := fmt.Sprintf("username:%s", username)
	userID, err := r.db.Get(ctx, usernameKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	userKey := fmt.Sprintf("user:%s", userID)
	result, err := r.db.HGetAll(ctx, userKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return redisMapToUser(result)
}

func redisMapToUser(m map[string]string) (*user.User, error) {
	var u user.User
	u.ID = m["id"]
	u.Username = m["username"]
	u.Password = m["password"]

	t, err := time.Parse(time.RFC3339, m["created_at"])
	if err != nil {
		return nil, err
	}
	u.CreatedAt = t

	return &u, nil
}
