package api

import (
	"chat-service/internal/config"
	"chat-service/internal/middleware"
	"chat-service/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter())

	// 设置配置到上下文
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// 初始化WebSocket Hub
	websocket.StartHub()

	// 控制器实例
	authController := NewAuthController()
	userController := NewUserController()
	chatController := NewChatController()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "chat-service",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(&cfg.JWT))
		{
		// 用户相关
		users := protected.Group("/users")
		{
			users.GET("/profile", userController.GetProfile)
			users.PUT("/profile", userController.UpdateProfile)
			users.GET("/search", userController.SearchUsers)
			users.GET("/:id", userController.GetUserByID)
		}

			// 聊天相关
			rooms := protected.Group("/rooms")
			{
				rooms.GET("", chatController.GetRooms)
				rooms.GET("/unread", chatController.GetRoomsWithUnread)
				rooms.POST("", chatController.CreateRoom)
				rooms.GET("/:id", chatController.GetRoom)
				rooms.POST("/:id/join", chatController.JoinRoom)
				rooms.POST("/:id/leave", chatController.LeaveRoom)
				rooms.GET("/:id/messages", chatController.GetMessages)
				rooms.POST("/:id/read", chatController.MarkAsRead)
				rooms.GET("/:id/unread", chatController.GetUnreadCount)
				rooms.GET("/:id/members", chatController.GetRoomMembers)
				rooms.POST("/:id/members", chatController.AddMember)
			}

			// WebSocket连接
			protected.GET("/ws", websocket.HandleWebSocket)
		}
	}

	return r
}
