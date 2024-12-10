package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper, log *logrus.Logger) *redis.Client {
	addr := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
	password := viper.GetString("redis.password")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	ctx := context.TODO()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("✅ Redis client connected successfully...")

	return rdb
}
