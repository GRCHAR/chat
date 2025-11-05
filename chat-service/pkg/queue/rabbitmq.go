package queue

import (
	"chat-service/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

type Message struct {
	Type      string      `json:"type"` // message, system, notification
	RoomID    uint        `json:"room_id"`
	SenderID  uint        `json:"sender_id"`
	Content   interface{} `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

func InitRabbitMQ(cfg *config.RabbitMQConfig) error {
	var err error

	// 连接到RabbitMQ
	conn, err = amqp.Dial(cfg.URL)
	if err != nil {
		return fmt.Errorf("RabbitMQ连接失败: %v", err)
	}

	// 创建通道
	channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("RabbitMQ通道创建失败: %v", err)
	}

	// 声明交换机
	err = channel.ExchangeDeclare(
		cfg.Exchange, // 交换机名称
		"direct",     // 交换机类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部使用
		false,        // 不等待
		nil,          // 参数
	)
	if err != nil {
		return fmt.Errorf("交换机声明失败: %v", err)
	}

	// 声明队列
	_, err = channel.QueueDeclare(
		cfg.Queue, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他
		false,     // 不等待
		nil,       // 参数
	)
	if err != nil {
		return fmt.Errorf("队列声明失败: %v", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		cfg.Queue,      // 队列名称
		"chat.message", // 路由键
		cfg.Exchange,   // 交换机名称
		false,          // 不等待
		nil,            // 参数
	)
	if err != nil {
		return fmt.Errorf("队列绑定失败: %v", err)
	}

	fmt.Println("RabbitMQ连接成功")
	return nil
}

// PublishMessage 发布消息
func PublishMessage(ctx context.Context, msgType string, roomID, senderID uint, content interface{}) error {
	message := Message{
		Type:      msgType,
		RoomID:    roomID,
		SenderID:  senderID,
		Content:   content,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return channel.Publish(
		"chat.exchange", // 交换机
		"chat.message",  // 路由键
		false,           // 强制
		false,           // 立即
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  "30000", // 30秒过期
		},
	)
}

// ConsumeMessages 消费消息
func ConsumeMessages(handler func(Message) error) error {
	msgs, err := channel.Consume(
		"chat.queue", // 队列
		"",           // 消费者标签
		false,        // 自动确认
		false,        // 排他
		false,        // 不等待
		false,        // 参数
		nil,
	)
	if err != nil {
		return fmt.Errorf("消费者注册失败: %v", err)
	}

	go func() {
		for msg := range msgs {
			var message Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("消息解析失败: %v", err)
				msg.Ack(false)
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("消息处理失败: %v", err)
				msg.Nack(false, true) // 重新入队
			} else {
				msg.Ack(false) // 确认消息
			}
		}
	}()

	return nil
}

// Close 关闭连接
func Close() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
