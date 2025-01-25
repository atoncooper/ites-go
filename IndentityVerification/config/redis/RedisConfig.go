package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	RedisConn *redis.Client
	SshClient *ssh.Client
)

type Config struct {
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	System struct {
		Auth   string `yaml:"auth"`
		System string `yaml:"password"`
	} `yaml:"system"`
}

func InitRedis() {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	yamlPath := filepath.Join(dir, "config", "redis", "redis.yaml")
	// 读取yaml文件
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		// 日志记录
		fmt.Println("reading yaml error")
		return
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// 日志记录
		fmt.Println("unmarshal yaml error")
		return
	}
	/*
		通过ssh连接redis
	*/
	//sshConfig := &ssh.ClientConfig{
	//	User:            config.System.Auth,
	//	Auth:            []ssh.AuthMethod{ssh.Password(config.System.System)},
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//	Timeout:         15 * time.Second,
	//}
	//remoteIp := fmt.Sprintf("%s:%s", config.Redis.Host, "22")
	//SshClient, err = ssh.Dial("tcp", remoteIp, sshConfig)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("ssh连接成功", SshClient)
	//RedisConn = redis.NewClient(&redis.Options{
	//	Addr: net.JoinHostPort(config.Redis.Host, config.Redis.Port),
	//	Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
	//		return SshClient.Dial(network, addr)
	//	},
	//	// Disable timeouts, because SSH does not support deadlines.
	//	ReadTimeout:  -1,
	//	WriteTimeout: -1,
	//})
	/*
		简单连接
	*/
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
}
