package database

import "github.com/redis/go-redis/v9"

type UserRepo struct {
	db *redis.Client
}

func NewUserRepo(db *redis.Client) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser()     {}
func (r *UserRepo) GetUserByEmail() {}
