package database

import (
	"chatter/server/internal/user"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("key not found")

type UserRepo struct {
	db *redis.Client
}

func NewUserRepo(db *redis.Client) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser(u *user.User) error {
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now().UTC()

	userKey := fmt.Sprintf("user:%s", u.ID)
	userNameKey := fmt.Sprintf("username:%s", u.Username)

	pipe := r.db.TxPipeline()

	set := pipe.SetNX(ctx, userNameKey, u.ID, 0)

	userJSON, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("databse: failed to marshall user: %v", err)
	}
	pipe.HSet(ctx, userKey, "data", userJSON)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("database: failed to complete redis transaction, %v", err)
	}

	if !set.Val() {
		return user.ErrUserAlreadyExists
	}

	return nil
}

func (r *UserRepo) GetUserByUsername(username string) (*user.User, error) {
	usernameKey := fmt.Sprintf("username:%s", username)
	userID, err := r.db.Get(ctx, usernameKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	userKey := fmt.Sprintf("user:%s", userID)
	userData, err := r.db.HGet(ctx, userKey, "data").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var u user.User
	if err := json.Unmarshal([]byte(userData), &u); err != nil {
		return nil, fmt.Errorf("database: failed to unmarshall user data, %v", err)
	}

	return &u, nil
}
