package main

import (
	pb "IndentityVerification/proto/CreateToken"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 默认不加密
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 建立建立
	client := pb.NewCreateLongTokenClient(conn)
	resp, _ := client.CreateLongToken(context.Background(), &pb.CreateLongTokenRequest{Username: "admin", Key: "123456"})
	fmt.Println(resp.GetLongToken())
}
