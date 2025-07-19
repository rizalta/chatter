package database

import (
	"chatter/server/internal/user"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("key not found")

type UserRepo struct {
	db *redis.Client
}

func NewUserRepo(db *redis.Client) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser(username, password string) (*user.User, error) {
	return nil, nil
}

func (r *UserRepo) GetUserByUsername(username string) (*user.User, error) {
	return nil, nil
}
