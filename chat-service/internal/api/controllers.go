package api

import (
	"chat-service/internal/config"
	"chat-service/internal/middleware"
	"chat-service/internal/models"
	"chat-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userService *service.UserService
}

type ChatController struct {
	chatService    *service.ChatService
	messageService *service.MessageService
}

type UserController struct {
	userService *service.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		userService: service.NewUserService(),
	}
}

func NewChatController() *ChatController {
	return &ChatController{
		chatService:    service.NewChatService(),
		messageService: service.NewMessageService(),
	}
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Nickname string `json:"nickname" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 创建房间请求结构
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	Type        string `json:"type" binding:"required,oneof=single group"`
	MemberIDs   []uint `json:"member_ids" binding:"required,min=1"`
}

// 发送消息请求结构
type SendMessageRequest struct {
	RoomID  uint   `json:"room_id" binding:"required"`
	Content string `json:"content" binding:"required,max=1000"`
	Type    string `json:"type" binding:"required,oneof=text image file"`
}

// 认证相关接口
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	if _, err := c.userService.GetUserByUsername(req.Username); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   "active",
	}

	if err := c.userService.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	// 生成JWT token
	cfg := ctx.MustGet("config").(*config.Config)
	token, err := middleware.GenerateToken(user.ID, &cfg.JWT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.GetUserByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	cfg := ctx.MustGet("config").(*config.Config)
	token, err := middleware.GenerateToken(user.ID, &cfg.JWT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// 用户相关接口
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req models.User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新允许修改的字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := c.userService.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// 搜索用户
func (c *UserController) SearchUsers(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	users, err := c.userService.SearchUsers(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "搜索用户失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

// 根据ID获取用户
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := c.userService.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// 聊天相关接口
func (c *ChatController) GetRooms(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	rooms, err := c.chatService.GetUserRooms(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取房间列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (c *ChatController) GetRoomsWithUnread(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	rooms, err := c.messageService.GetRoomsWithUnread(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取房间列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (c *ChatController) CreateRoom(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req CreateRoomRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := &models.ChatRoom{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		OwnerID:     userID,
		MaxMembers:  100,
	}

	if req.Type == "single" && len(req.MemberIDs) == 1 {
		// 单聊，将对方也加入成员列表
		req.MemberIDs = append(req.MemberIDs, userID)
	} else {
		// 群聊，确保房主在成员列表中
		hasOwner := false
		for _, id := range req.MemberIDs {
			if id == userID {
				hasOwner = true
				break
			}
		}
		if !hasOwner {
			req.MemberIDs = append(req.MemberIDs, userID)
		}
	}

	if err := c.chatService.CreateRoom(room, req.MemberIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "房间创建失败"})
		return
	}

	room, err := c.chatService.GetRoomByID(room.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取房间信息失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"room": room})
}

func (c *ChatController) GetRoom(ctx *gin.Context) {
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	room, err := c.chatService.GetRoomByID(uint(roomID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"room": room})
}

func (c *ChatController) JoinRoom(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	if err := c.chatService.JoinRoom(userID, uint(roomID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "加入房间失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "成功加入房间"})
}

func (c *ChatController) LeaveRoom(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	if err := c.chatService.LeaveRoom(userID, uint(roomID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "离开房间失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "成功离开房间"})
}

func (c *ChatController) GetMessages(ctx *gin.Context) {
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	messages, err := c.messageService.GetRoomMessages(uint(roomID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages":  messages,
		"page":      page,
		"page_size": pageSize,
	})
}

func (c *ChatController) MarkAsRead(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	if err := c.messageService.MarkAsRead(userID, uint(roomID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "标记已读失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}

func (c *ChatController) GetUnreadCount(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	count, err := c.messageService.GetUnreadCount(userID, uint(roomID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取未读数失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (c *ChatController) GetRoomMembers(ctx *gin.Context) {
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	members, err := c.chatService.GetRoomMembers(uint(roomID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"members": members})
}

// 添加成员到房间
type AddMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

func (c *ChatController) AddMember(ctx *gin.Context) {
	roomID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的房间ID"})
		return
	}

	var req AddMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查房间是否存在
	room, err := c.chatService.GetRoomByID(uint(roomID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	// 只有群聊才能添加成员
	if room.Type != "group" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "只有群聊才能添加成员"})
		return
	}

	// 添加成员
	if err := c.chatService.JoinRoom(req.UserID, uint(roomID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "添加成员失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "成员添加成功"})
}
