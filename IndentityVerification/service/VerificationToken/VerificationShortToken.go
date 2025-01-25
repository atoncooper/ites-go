package VerificationToken

import (
	"IndentityVerification/models/redis"
	pb "IndentityVerification/proto/VerificationToken"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type VerificationShortTokenServer struct {
	pb.UnimplementedVerificationShortTokenServer
}

func (v *VerificationShortTokenServer) VerificationShortToken(ctx context.Context, req *pb.VerificationShortTokenRequest) (*pb.VerificationShortTokenResponse, error) {
	shortToken := req.ShortToken
	obj := redis.UserInfo{UserID: req.Username}
	shortPrivateKeyStr, err := obj.GetPrivateKey()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	shortPrivateKey, err := base64.StdEncoding.DecodeString(shortPrivateKeyStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	token, err := jwt.ParseWithClaims(shortToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return shortPrivateKey, nil
	})
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, fmt.Errorf("令牌已过期")
		}
		return &pb.VerificationShortTokenResponse{IsValid: true}, nil
	} else {
		return &pb.VerificationShortTokenResponse{IsValid: false}, nil
	}
}
