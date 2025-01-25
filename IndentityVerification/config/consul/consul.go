package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
)

// ConsulConfig 用于解析Consul配置
type ConsulConfig struct {
	Address string `yaml:"address"`
	Scheme  string `yaml:"scheme"`
	// Datacenter string `yaml:"datacenter"`
	Token string `yaml:"token"`
}

// GrpcConfig 用于解析gRPC配置
type GrpcConfig struct {
	ServiceName     string `yaml:"serviceName"`
	Address         string `yaml:"address"`
	Port            int    `yaml:"port"`
	HealthCheckPort int    `yaml:"healthCheckPort"`
}

// Config 用于解析整个配置文件
type Config struct {
	Consul ConsulConfig `yaml:"consul"`
	Grpc   GrpcConfig   `yaml:"grpc"`
}

func loadConfig(filePath string) (*Config, error) {
	// ioutil.ReadFile(filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

func registerServiceWithConsul(cfg *api.Config, serviceName, serviceAddress string, port int) error {
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	registration := &api.AgentServiceRegistration{
		Name:    serviceName,
		Address: serviceAddress,
		Port:    port,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/grpc.health.v1.Health/Check", serviceAddress, port),
			Interval:                       "10s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "5m",
		},
	}
	d := fmt.Sprintf("%s:%d/grpc.health.v1.Health/Check", serviceAddress, port)
	fmt.Println(d)
	return client.Agent().ServiceRegister(registration)
}

func main() {
	// 加载配置文件
	dir, err := filepath.Abs(filepath.Dir("."))
	yamlPath := filepath.Join(dir, "config", "consul", "consul.yaml")
	config, err := loadConfig(yamlPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 初始化Consul客户端配置
	consulCfg := api.DefaultConfig()
	consulCfg.Address = config.Consul.Address
	consulCfg.Scheme = config.Consul.Scheme
	// consulCfg.Datacenter = config.Consul.Datacenter
	consulCfg.Token = config.Consul.Token

	// 注册服务到Consul
	serviceName := config.Grpc.ServiceName
	serviceAddress := config.Grpc.Address
	// servicePort := config.Grpc.Port
	healthCheckPort := config.Grpc.HealthCheckPort

	if err := registerServiceWithConsul(consulCfg, serviceName, serviceAddress, healthCheckPort); err != nil {
		log.Fatalf("failed to register service with Consul: %v", err)
	}

	log.Printf("Service %s registered with Consul", serviceName)

	// 等待信号终止
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// 注销服务
	client, _ := api.NewClient(consulCfg)
	err = client.Agent().ServiceDeregister(serviceName)
	if err != nil {
		return
	}
	log.Printf("Service %s deregistered", serviceName)
}
