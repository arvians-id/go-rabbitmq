package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func InitRedis(configuration Config) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	redisHost := configuration.Get("REDIS_HOST")
	redisPort := configuration.Get("REDIS_PORT")
	redisPassword := configuration.Get("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(configuration.Get("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       redisDB,
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
