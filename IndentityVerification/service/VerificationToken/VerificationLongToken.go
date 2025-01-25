package VerificationToken

import (
	"IndentityVerification/models/redis"
	pb "IndentityVerification/proto/VerificationToken"
	"IndentityVerification/service/CreateToken"
	"IndentityVerification/service/RabbitMQ"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"log"
	"time"
)

type VerificationLongTokenServer struct {
	pb.UnimplementedVerificationLongTokenServer
}

func HandleAmqpCall(q *CreateToken.QueueInfo) {
	/*
		处理rabbitmq服务的调用
	*/
	msgs, err := RabbitMQ.ConsumeMessages(q.SendQueue)
	if err != nil {
		// TODO 进行日志记录
	}
	func() {
		for msg := range msgs {
			var message map[string]interface{}
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			req := &pb.VerificationLongTokenRequest{Token: message["token"].(string), Username: message["username"].(string)}
			conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
			if err != nil {
				fmt.Printf("无法连接到gRPC服务器: %v\n", err)
				return
			}
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {

				}
			}(conn)
			client := pb.NewVerificationLongTokenClient(conn)
			resp, _ := client.VerificationLongToken(context.Background(), req)
			// 返回处理结果给另一个队列
			err = RabbitMQ.PublishMessage(q.ReceiveQueue, map[string]interface{}{"result": resp.GetIsValid()})
			if err != nil {
				return
			}
		}
	}()
}
func (v *VerificationLongTokenServer) VerificationLongToken(ctx context.Context, req *pb.VerificationLongTokenRequest) (*pb.VerificationLongTokenResponse, error) {
	longToken := req.Token
	obj := redis.UserInfo{UserID: req.Username}
	privateKeyStr, err := obj.GetLongPrivateKey()
	if err != nil {
		return nil, err
	}
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		// TODO 进行日志记录
	}
	token, err := jwt.ParseWithClaims(longToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名算法")
		}
		return privateKey, nil
	})
	// 检查令牌是否有效（这里验证了签名等基础信息后，再验证有效期等其他声明）
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// 检查是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, fmt.Errorf("令牌已过期")
		}
		// 这里还可以添加更多验证逻辑，比如验证Issuer、Subject等是否符合预期

		// 如果一切验证通过，构造并返回成功的响应
		return &pb.VerificationLongTokenResponse{
			IsValid: true,
		}, nil
	}
	return &pb.VerificationLongTokenResponse{
		IsValid: false,
	}, nil
}
