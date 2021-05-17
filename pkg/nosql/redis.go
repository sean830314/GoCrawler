package nosql

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

type RedisClient struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisDB() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	return redisClient, nil
}
