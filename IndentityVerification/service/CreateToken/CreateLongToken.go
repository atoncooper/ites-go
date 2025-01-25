package CreateToken

import (
	"IndentityVerification/models/redis"
	pb "IndentityVerification/proto/CreateToken"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CreateLongTokenServer struct {
	pb.UnimplementedCreateLongTokenServer
}

func (c *CreateLongTokenServer) CreateLongToken(ctx context.Context, req *pb.CreateLongTokenRequest) (*pb.CreateLongTokenResponse, error) {
	longPrivateKey := CreatePrivateKey()
	claims := jwt.RegisteredClaims{
		Issuer:    "ites.s.com",
		Subject:   req.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour * 24 * 7)), // 设置7天过期
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        req.Key,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名令牌
	tokenString, err := token.SignedString(longPrivateKey)
	if err != nil {
		// TODO 进行日志记录
		fmt.Println(err)
		return nil, err
	}
	// 对密钥进行存储
	save := redis.UserInfo{
		UserID:     req.Username,
		PrivateKey: base64.StdEncoding.EncodeToString(longPrivateKey),
	}
	err = save.SaveLongPrivateKey()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 返回token
	return &pb.CreateLongTokenResponse{LongToken: tokenString}, nil
}
