package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	Conn *amqp.Connection
	Ch   *amqp.Channel
)

type Config struct {
	RabbitMQ struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"rabbitmq"`
	System struct {
		Auth   string `yaml:"auth"`
		System string `yaml:"password"`
	} `yaml:"system"`
}

func InitRabbitMQ() { // 初始化RabbitMQ
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	yamlPath := filepath.Join(dir, "config", "rabbitmq", "rabbitmq.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	// 连接RabbitMQ
	Conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.System.Auth, config.System.System, config.RabbitMQ.Host, config.RabbitMQ.Port))
	if err != nil {
		// TODO 进行日志记录
		panic(err)
	}
	defer Conn.Close()
	Ch, err = Conn.Channel()
	defer Ch.Close()
}
