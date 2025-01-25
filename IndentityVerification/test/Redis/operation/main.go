package main

import (
	"IndentityVerification/config/redis"
	"context"
	"fmt"
	"time"
)

func main() {
	redis.InitRedis() // 实例化
	ctx := context.Background()
	redis.RedisConn.Set(ctx, "key", "value2", time.Second*600)
	s := redis.RedisConn.Get(ctx, "key")
	fmt.Println(s)
	defer redis.RedisConn.Close()
	// redis.CloseRedis()
}
