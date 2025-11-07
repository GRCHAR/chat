package service

import (
	"chat-service/internal/database"
	"chat-service/internal/models"
	"chat-service/pkg/cache"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserService struct{}

type ChatService struct{}

type MessageService struct{}

// 用户相关服务
func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user *models.User) error {
	return database.GetDB().Create(user).Error
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.GetDB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	return database.GetDB().Save(user).Error
}

func (s *UserService) SearchUsers(query string) ([]models.User, error) {
	var users []models.User
	err := database.GetDB().
		Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", 
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(20).
		Find(&users).Error
	return users, err
}

// 聊天室相关服务
func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) CreateRoom(room *models.ChatRoom, memberIDs []uint) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 创建聊天室
		if err := tx.Create(room).Error; err != nil {
			return err
		}

		// 添加所有成员
		for _, memberID := range memberIDs {
			member := &models.RoomMember{
				RoomID:   room.ID,
				UserID:   memberID,
				Role:     "member",
				JoinedAt: time.Now(),
			}
			if memberID == room.OwnerID {
				member.Role = "owner"
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ChatService) GetUserRooms(userID uint) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := database.GetDB().
		Joins("JOIN room_members ON chat_rooms.id = room_members.room_id").
		Where("room_members.user_id = ? AND room_members.deleted_at IS NULL", userID).
		Preload("Owner").
		Find(&rooms).Error
	return rooms, err
}

func (s *ChatService) GetRoomByID(id uint) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := database.GetDB().
		Preload("Owner").
		Preload("Members.User").
		First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *ChatService) JoinRoom(userID, roomID uint) error {
	// 检查用户是否已经是成员
	var count int64
	database.GetDB().Model(&models.RoomMember{}).
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Count(&count)

	if count > 0 {
		return nil // 已经是成员
	}

	member := &models.RoomMember{
		RoomID:   roomID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	return database.GetDB().Create(member).Error
}

func (s *ChatService) LeaveRoom(userID, roomID uint) error {
	return database.GetDB().Where("user_id = ? AND room_id = ?", userID, roomID).
		Delete(&models.RoomMember{}).Error
}

func (s *ChatService) GetRoomMembers(roomID uint) ([]models.RoomMember, error) {
	var members []models.RoomMember
	err := database.GetDB().
		Preload("User").
		Where("room_id = ?", roomID).
		Find(&members).Error
	return members, err
}

// 消息相关服务
func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) CreateMessage(message *models.Message) error {
	return database.GetDB().Create(message).Error
}

func (s *MessageService) GetRoomMessages(roomID uint, page, pageSize int) ([]models.Message, error) {
	var messages []models.Message
	offset := (page - 1) * pageSize

	err := database.GetDB().
		Preload("Sender").
		Where("room_id = ? AND is_deleted = false", roomID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}

func (s *MessageService) GetMessageByID(id uint) (*models.Message, error) {
	var message models.Message
	err := database.GetDB().
		Preload("Sender").
		Preload("Room").
		First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *MessageService) DeleteMessage(id uint) error {
	return database.GetDB().Model(&models.Message{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (s *MessageService) GetUnreadCount(userID, roomID uint) (int64, error) {
	var count int64

	// 先获取用户在该房间的最后阅读时间
	var lastRead time.Time
	database.GetDB().Model(&models.RoomMember{}).
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Pluck("last_read", &lastRead)

	// 统计未读消息数
	database.GetDB().Model(&models.Message{}).
		Where("room_id = ? AND sender_id != ? AND is_deleted = false AND created_at > ?",
			roomID, userID, lastRead).
		Count(&count)

	return count, nil
}

func (s *MessageService) MarkAsRead(userID, roomID uint) error {
	return database.GetDB().Model(&models.RoomMember{}).
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Update("last_read", time.Now()).Error
}

func (s *MessageService) GetRoomsWithUnread(userID uint) ([]map[string]interface{}, error) {
	type RoomWithUnread struct {
		models.ChatRoom
		UnreadCount int64           `json:"unread_count"`
		LastMessage *models.Message `json:"last_message,omitempty"`
	}

	var results []RoomWithUnread

	err := database.GetDB().
		Table("chat_rooms").
		Select(`
			chat_rooms.*, 
			(
				SELECT COUNT(*) 
				FROM messages 
				WHERE messages.room_id = chat_rooms.id 
				AND messages.sender_id != ? 
				AND messages.is_deleted = false 
				AND messages.created_at > COALESCE(
					(SELECT last_read FROM room_members WHERE user_id = ? AND room_id = chat_rooms.id),
					'1970-01-01'
				)
			) as unread_count,
			(
				SELECT * 
				FROM messages 
				WHERE messages.room_id = chat_rooms.id 
				AND messages.is_deleted = false 
				ORDER BY created_at DESC 
				LIMIT 1
			) as last_message
		`, userID, userID).
		Joins("JOIN room_members ON chat_rooms.id = room_members.room_id").
		Where("room_members.user_id = ?", userID).
		Preload("LastMessage.Sender").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为map格式
	var rooms []map[string]interface{}
	for _, result := range results {
		room := map[string]interface{}{
			"id":           result.ID,
			"name":         result.Name,
			"description":  result.Description,
			"type":         result.Type,
			"avatar":       result.Avatar,
			"owner_id":     result.OwnerID,
			"max_members":  result.MaxMembers,
			"created_at":   result.CreatedAt,
			"updated_at":   result.UpdatedAt,
			"unread_count": result.UnreadCount,
			"last_message": result.LastMessage,
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// 获取在线用户数
func (s *ChatService) GetOnlineUserCount(roomID uint) (int, error) {
	userIDs, err := cache.GetRoomUsers(context.Background(), roomID)
	if err != nil {
		return 0, err
	}
	return len(userIDs), nil
}

// 获取用户是否在线
func (s *UserService) IsUserOnline(userID uint) (bool, error) {
	_, err := cache.GetUserOnline(context.Background(), userID)
	if err != nil {
		return false, nil
	}
	return true, nil
}
