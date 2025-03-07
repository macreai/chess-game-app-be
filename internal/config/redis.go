package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_ADDRESS"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		Protocol: viper.GetInt("REDIS_PROTOCOL"),
	})
}
