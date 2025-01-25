package redis

import (
	"IndentityVerification/config/redis"
	"context"
	"fmt"
	"time"
)

type UserInfo struct {
	UserID     string
	PrivateKey string
}

func (u UserInfo) SavaPrivateKey() error {
	//@description 将短token密钥存储进redis
	ctx := context.Background()
	err := redis.RedisConn.Set(ctx, fmt.Sprintf("%s:PrivateKey", u.UserID), u.PrivateKey, time.Second*10).Err()
	if err != nil {
		return err
	}
	redis.RedisConn.Get(ctx, fmt.Sprintf("%s:PrivateKey", u.UserID))
	return nil
}
func (u UserInfo) GetPrivateKey() (string, error) {
	/*
		@description 获取短token密钥
	*/
	shortPrivateKey, err := redis.RedisConn.Get(context.Background(), fmt.Sprintf("%s:PrivateKey", u.UserID)).Result()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return shortPrivateKey, nil
}

func (u UserInfo) SaveLongPrivateKey() error {
	// @description 将长token密钥存储进redis
	ctx := context.Background()
	err := redis.RedisConn.Set(ctx, fmt.Sprintf("%s:LongPrivateKey", u.UserID), u.PrivateKey, time.Hour*24*7).Err() // 设置7天过期
	if err != nil {
		// TODO 进行日志记录
		return err
	}
	return nil
}
func (u UserInfo) GetLongPrivateKey() (string, error) {
	privateKey, err := redis.RedisConn.Get(context.Background(), fmt.Sprintf("%s:LongPrivateKey", u.UserID)).Result()
	if err != nil {
		// TODO 进行日志记录
		return "", err
	}
	return privateKey, err
}
