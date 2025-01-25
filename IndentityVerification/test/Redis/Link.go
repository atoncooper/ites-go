package main

import (
	"IndentityVerification/config/redis"
	"context"
)

func main() {
	ctx := context.Background()
	redis.RedisConn.Set(ctx, "key", "value", 100)
}
