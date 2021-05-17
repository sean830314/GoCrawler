package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sean830314/GoCrawler/pkg/nosql"
	"github.com/sirupsen/logrus"
)

type AuthInterface interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteToken(string) error
	DeleteAuth(*AccessDetails) error
}

type RedisAuthService struct {
	client *redis.Client
}

var _ AuthInterface = &RedisAuthService{}

func NewRedisAuthService() *RedisAuthService {
	redisClient, err := nosql.NewRedisDB()
	if err != nil {
		logrus.Error("New NewRedisDB client failed, error: ", err)
	}
	return &RedisAuthService{client: redisClient}
}

//Save token metadata to Redis
func (rs *RedisAuthService) CreateAuth(userId string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	if err := rs.client.Set(td.TokenUuid, userId, at.Sub(now)).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	if err := rs.client.Set(td.RefreshUuid, userId, rt.Sub(now)).Err(); err != nil {
		return err
	}
	return nil
}

//Get userID by tokenUuid in redis
func (rs *RedisAuthService) FetchAuth(tokenUuid string) (string, error) {
	userid, err := rs.client.Get(tokenUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

//Once a user row in the token table
func (rs *RedisAuthService) DeleteAuth(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	//delete access token
	if err := rs.DeleteToken(authD.TokenUuid); err != nil {
		return err
	}
	//delete refresh token
	if err := rs.DeleteToken(refreshUuid); err != nil {
		return err
	}
	return nil
}

func (rs *RedisAuthService) DeleteToken(token string) error {
	//delete token
	deleted, err := rs.client.Del(token).Result()
	if err != nil {
		return err
	}
	if deleted != 1 {
		return errors.New("delete refresh token error in redis")
	}
	return nil
}
