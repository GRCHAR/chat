package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50" json:"username"`
	Nickname  string         `gorm:"size:100" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Email     string         `gorm:"uniqueIndex;size:100" json:"email"`
	Password  string         `gorm:"size:255" json:"-"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Type        string         `gorm:"size:20;default:'group'" json:"type"` // single, group
	Avatar      string         `gorm:"size:255" json:"avatar"`
	OwnerID     uint           `json:"owner_id"`
	MaxMembers  int            `gorm:"default:100" json:"max_members"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Owner    User         `gorm:"foreignKey:OwnerID" json:"owner"`
	Members  []RoomMember `gorm:"foreignKey:RoomID" json:"members"`
	Messages []Message    `gorm:"foreignKey:RoomID" json:"messages,omitempty"`
}

// RoomMember 聊天室成员
type RoomMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoomID    uint           `json:"room_id"`
	UserID    uint           `json:"user_id"`
	Role      string         `gorm:"size:20;default:'member'" json:"role"` // owner, admin, member
	JoinedAt  time.Time      `json:"joined_at"`
	LastRead  time.Time      `json:"last_read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Room ChatRoom `gorm:"foreignKey:RoomID" json:"room"`
	User User     `gorm:"foreignKey:UserID" json:"user"`
}

// Message 消息模型
type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoomID    uint           `json:"room_id"`
	SenderID  uint           `json:"sender_id"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:20;default:'text'" json:"type"` // text, image, file, system
	ReplyToID *uint          `json:"reply_to_id"`                        // 回复的消息ID
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Sender  User     `gorm:"foreignKey:SenderID" json:"sender"`
	Room    ChatRoom `gorm:"foreignKey:RoomID" json:"room"`
	ReplyTo *Message `gorm:"foreignKey:ReplyToID" json:"reply_to_message,omitempty"`
}

// UnreadMessage 未读消息计数
type UnreadMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	RoomID    uint      `json:"room_id"`
	Count     int       `gorm:"default:0" json:"count"`
	LastMsgID uint      `json:"last_msg_id"`
	UpdatedAt time.Time `json:"updated_at"`

	User User     `gorm:"foreignKey:UserID" json:"user"`
	Room ChatRoom `gorm:"foreignKey:RoomID" json:"room"`
}

// OnlineUser 在线用户
type OnlineUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	ConnID    string    `gorm:"size:255;uniqueIndex" json:"conn_id"`
	RoomIDs   string    `gorm:"type:text" json:"room_ids"` // JSON数组
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
