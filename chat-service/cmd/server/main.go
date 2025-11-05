package main

import (
	"chat-service/internal/api"
	"chat-service/internal/config"
	"chat-service/internal/database"
	"chat-service/pkg/cache"
	"chat-service/pkg/queue"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化Redis
	if err := cache.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}

	// 初始化RabbitMQ
	if err := queue.InitRabbitMQ(&cfg.RabbitMQ); err != nil {
		log.Printf("RabbitMQ初始化失败: %v", err)
		log.Println("系统将在没有消息队列的情况下运行")
	} else {
		// 启动消息队列消费者
		go startMessageQueueConsumer()
	}

	// 设置路由
	router := api.SetupRouter(cfg)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("聊天服务启动在 %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器强制关闭: %v", err)
	}

	// 关闭RabbitMQ连接
	queue.Close()

	log.Println("服务器已关闭")
}

// 启动消息队列消费者
func startMessageQueueConsumer() {
	err := queue.ConsumeMessages(func(msg queue.Message) error {
		log.Printf("收到消息: %+v", msg)

		// 在这里处理消息队列中的消息
		// 例如：发送推送通知、记录日志、处理离线消息等

		return nil
	})

	if err != nil {
		log.Printf("消息队列消费者启动失败: %v", err)
	}
}
