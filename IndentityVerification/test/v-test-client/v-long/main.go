package main

import (
	pb "IndentityVerification/proto/VerificationToken"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	// 连接到gRPC服务器，根据实际情况调整地址和端口
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("无法连接到gRPC服务器: %v\n", err)
		return
	}
	defer conn.Close()

	// 创建验证长令牌服务的客户端实例
	client := pb.NewVerificationLongTokenClient(conn)
	resp, _ := client.VerificationLongToken(context.Background(), &pb.VerificationLongTokenRequest{Username: "admin", Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpdGVzLnMuY29tIiwic3ViIjoiYWRtaW4iLCJleHAiOjE3Mzk1MzczMzAsIm5iZiI6MTczNjUxMzMzMCwianRpIjoiMTIzNDU2In0.g3t3InYOQeujxGiy5y806lha88vwCyN1FDwjIShccrk"})
	fmt.Println(resp.GetIsValid())
}
