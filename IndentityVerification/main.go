package main

import (
	"IndentityVerification/config/rabbitmq"
	"IndentityVerification/config/redis"
	pb "IndentityVerification/proto/CreateToken"
	vpb "IndentityVerification/proto/VerificationToken"
	"IndentityVerification/service/CreateToken"
	"IndentityVerification/service/VerificationToken"
	"IndentityVerification/utils/LogManager"
	"errors"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

func registerService(client *api.Client) {
	/*
		@description 注册服务到consul
	*/
	registration := &api.AgentServiceRegistration{
		ID:      "IdentityVerification-ID",
		Name:    "IdentityVerification-Server",
		Address: "127.0.0.1",
		Port:    9090,
		Check: &api.AgentServiceCheck{
			GRPC:                           "127.0.0.1:9090/grpc.health.v1.Health/Check",
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

}

func InitConsul() {
	/*
		@description 初始化里连接consul
	*/
	config := api.DefaultConfig()
	config.Address = "localhost:8500"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	// 服务注册
	registerService(client)
}
func main() {
	// 初始化日志
	logManager := LogManager.Info{}
	logManager.NewLogManager()
	defer logManager.CloseLogFiles()
	// 初始化rabbitmq
	rabbitmq.InitRabbitMQ()
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			logManager.ErrorLog("StartRabbiMQ", err)
		}
	}(rabbitmq.Conn)
	// 初始化redis连接
	redis.InitRedis()
	defer func(RedisConn *redis2.Client) {
		err := RedisConn.Close()
		if err != nil {

		}
	}(redis.RedisConn)
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		// TODO 进行日志记录
		logManager.ErrorLog("StartGrpc", err)
	}
	// 创建grpc服务
	conn := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("IdentityVerification-Server", grpc_health_v1.HealthCheckResponse_SERVING)
	pb.RegisterCreateTokenServer(conn, &CreateToken.CreateShortTokenServer{})                         // 短token服务
	pb.RegisterCreateLongTokenServer(conn, &CreateToken.CreateLongTokenServer{})                      // 长token服务
	vpb.RegisterVerificationLongTokenServer(conn, &VerificationToken.VerificationLongTokenServer{})   // 校验长token服务
	vpb.RegisterVerificationShortTokenServer(conn, &VerificationToken.VerificationShortTokenServer{}) // 校验短token服务
	grpc_health_v1.RegisterHealthServer(conn, healthServer)
	logManager.InfoLog("StartGrpc", errors.New("grpc服务启动成功"))
	fmt.Println("Server Start Success")
	err = conn.Serve(listen)

}
