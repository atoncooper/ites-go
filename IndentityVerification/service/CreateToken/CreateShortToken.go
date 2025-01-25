package CreateToken

import (
	"IndentityVerification/models/redis"
	pb "IndentityVerification/proto/CreateToken"
	vpb "IndentityVerification/proto/VerificationToken"
	"IndentityVerification/service/RabbitMQ"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"time"
)

type CreateShortTokenServer struct {
	pb.UnimplementedCreateTokenServer
}

func CreatePrivateKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// TODO 进行日志记录
		panic("Failed to generate random key")
	}

	// 将密钥编码为 Base64 字符串
	// keyBase64 := base64.StdEncoding.EncodeToString(key)
	return key
}

type QueueInfo struct {
	SendQueue    string
	ReceiveQueue string
}

func (q *QueueInfo) VLongToken(token string, userName string) bool {
	// 设置队列名
	q.SendQueue = fmt.Sprintf("%s:q.SendQueue", userName)     // 发送消息队列
	q.ReceiveQueue = fmt.Sprintf("%s:reciverQueue", userName) // 接收消息队列
	err := RabbitMQ.DeclareQueue(q.SendQueue)
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	err = RabbitMQ.DeclareQueue(q.ReceiveQueue)
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	// 发送消息到队列
	dataMap := map[string]interface{}{
		"token":    token,
		"username": userName,
	}
	err = RabbitMQ.PublishMessage(q.SendQueue, dataMap)
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	// 阻塞等待返回数据
	msgs, err := RabbitMQ.ConsumeMessages(q.ReceiveQueue)
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	done := make(chan bool, 1)
	go func() {
		for msg := range msgs {
			if string(msg.Body) == "true" {
				done <- true
				break
			}
		}
	}()
	select {
	case <-done:
		return true
	case <-time.After(10 * time.Second):
		return false
	}
}

func Vt(token string, userName string) bool {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("无法连接到gRPC服务器: %v\n", err)
		return false
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	client := vpb.NewVerificationLongTokenClient(conn)
	resp, _ := client.VerificationLongToken(context.Background(), &vpb.VerificationLongTokenRequest{Username: userName, Token: token})
	return resp.GetIsValid()
}

func (c *CreateShortTokenServer) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.CreateTokenResponse, error) {
	/*
		@description 进行生成短token服务
		@return *pb.CreateTokenResponse 返回短token
	*/
	// TODO 对长token进行验证
	isValid := Vt(req.LongToken, req.Username)
	if !isValid {
		return &pb.CreateTokenResponse{Token: ""}, errors.New("长token无效")
	}
	privateKey := CreatePrivateKey()
	claims := jwt.RegisteredClaims{
		Issuer:    "ites.s.com",
		Subject:   req.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        req.Key,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名令牌
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		// TODO 进行日志记录
		fmt.Println(err)
		return nil, err
	}
	// TODO 存储privateKey相关信息
	save := redis.UserInfo{
		UserID:     req.Username,
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
	}
	err = save.SavaPrivateKey()
	if err != nil {
		// TODO 进行日志记录
		fmt.Println(err)
		return nil, err
	}
	return &pb.CreateTokenResponse{Token: tokenString}, nil
}
