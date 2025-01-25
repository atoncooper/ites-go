package RabbitMQ

import (
	"IndentityVerification/config/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := rabbitmq.Ch.Consume(
		queueName, // 队列名称
		"",        // 消费者标签
		true,      // 是否自动确认
		false,     // 是否独占
		false,     // 是否阻塞
		false,     // 额外参数
		nil,       // 额外参数
	)
	if err != nil {
		log.Printf("Failed to start consuming messages from queue %s: %v", queueName, err)
		return nil, err
	}
	// 返回消息通道
	return msgs, nil
}
