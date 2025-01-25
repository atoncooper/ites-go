package RabbitMQ

import (
	"IndentityVerification/config/rabbitmq"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareQueue(queueName string) error {
	/*
		DeclareQueue 声明队列
	*/
	_, err := rabbitmq.Ch.QueueDeclare(
		queueName, // 队列名称
		false,     // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占
		false,     // 是否阻塞
		nil,       // 额外参数
	)
	return err
}
func PublishMessage(queueName string, message map[string]interface{}) error {
	/*
		生产者生产数据
	*/
	body, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	return rabbitmq.Ch.Publish(
		"",        // 交换机名称
		queueName, // 队列名称
		false,     // 是否强制
		false,     // 是否立即
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}
