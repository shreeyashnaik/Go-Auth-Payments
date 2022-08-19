package db

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var (
	RedisDB *redis.Client
	ctx     = context.Background()
)

func SetDB() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func CreateUser(email, name, password string) error {
	if RedisDB.Exists(ctx, email).Val() == 1 {
		return errors.New("user already exists with this email")
	}

	RedisDB.LPush(ctx, email, name, password)
	return nil
}
