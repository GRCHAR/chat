package websocket

import (
	"chat-service/internal/database"
	"chat-service/internal/models"
	"chat-service/pkg/cache"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要检查origin
	},
}

type Client struct {
	ID     uint
	ConnID string
	Conn   *websocket.Conn
	Send   chan []byte
	Rooms  map[uint]bool
	mu     sync.RWMutex
}

type Hub struct {
	clients    map[string]*Client          // connID -> client
	rooms      map[uint]map[string]*Client // roomID -> connID -> client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

var hub = &Hub{
	clients:    make(map[string]*Client),
	rooms:      make(map[uint]map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan []byte),
}

type WSMessage struct {
	Type     string      `json:"type"`
	RoomID   uint        `json:"room_id"`
	SenderID uint        `json:"sender_id"`
	Content  interface{} `json:"content"`
	Time     time.Time   `json:"time"`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ConnID] = client
			h.mu.Unlock()
			log.Printf("客户端注册: %s", client.ConnID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ConnID]; ok {
				delete(h.clients, client.ConnID)
				close(client.Send)

				// 从所有房间中移除
				for roomID := range client.Rooms {
					if room, ok := h.rooms[roomID]; ok {
						delete(room, client.ConnID)
						if len(room) == 0 {
							delete(h.rooms, roomID)
						}
					}

					// 从Redis中移除
					cache.RemoveUserFromRoom(context.Background(), roomID, client.ID)
				}

				// 设置用户离线
				cache.SetUserOffline(context.Background(), client.ID)
			}
			h.mu.Unlock()
			log.Printf("客户端注销: %s", client.ConnID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ConnID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[string]*Client)
	}

	h.rooms[roomID][client.ConnID] = client
	client.Rooms[roomID] = true

	// 添加到Redis
	cache.AddUserToRoom(context.Background(), roomID, client.ID)

	log.Printf("用户 %d 加入房间 %d", client.ID, roomID)
}

func (h *Hub) LeaveRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		delete(room, client.ConnID)
		if len(room) == 0 {
			delete(h.rooms, roomID)
		}
	}

	delete(client.Rooms, roomID)

	// 从Redis中移除
	cache.RemoveUserFromRoom(context.Background(), roomID, client.ID)

	log.Printf("用户 %d 离开房间 %d", client.ID, roomID)
}

func (h *Hub) BroadcastToRoom(roomID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		for _, client := range room {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client.ConnID)
				delete(room, client.ConnID)
			}
		}
	}
}

func HandleWebSocket(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	connID := generateConnID()
	client := &Client{
		ID:     userID,
		ConnID: connID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Rooms:  make(map[uint]bool),
	}

	hub.register <- client

	// 设置用户在线
	roomIDs := getUserRoomIDs(userID)
	cache.SetUserOnline(context.Background(), userID, connID, roomIDs)

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var wsMsg WSMessage
		err := c.Conn.ReadJSON(&wsMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		wsMsg.SenderID = c.ID
		wsMsg.Time = time.Now()

		switch wsMsg.Type {
		case "join_room":
			roomID := wsMsg.RoomID
			if isValidRoomMember(c.ID, roomID) {
				hub.JoinRoom(c, roomID)
				c.SendMessage(WSMessage{
					Type:    "room_joined",
					RoomID:  roomID,
					Content: "成功加入房间",
					Time:    time.Now(),
				})
			} else {
				c.SendMessage(WSMessage{
					Type:    "error",
					Content: "无权限加入该房间",
					Time:    time.Now(),
				})
			}

		case "leave_room":
			hub.LeaveRoom(c, wsMsg.RoomID)
			c.SendMessage(WSMessage{
				Type:    "room_left",
				RoomID:  wsMsg.RoomID,
				Content: "成功离开房间",
				Time:    time.Now(),
			})

		case "message":
			roomID := wsMsg.RoomID
			if c.Rooms[roomID] {
				// 保存消息到数据库
				msg := &models.Message{
					RoomID:   roomID,
					SenderID: c.ID,
					Content:  wsMsg.Content.(string),
					Type:     "text",
				}

				if err := database.GetDB().Create(msg).Error; err != nil {
					log.Printf("消息保存失败: %v", err)
					continue
				}

				// 广播消息到房间
				messageData := WSMessage{
					Type:     "new_message",
					RoomID:   roomID,
					SenderID: c.ID,
					Content:  wsMsg.Content,
					Time:     time.Now(),
				}

				data, _ := json.Marshal(messageData)
				hub.BroadcastToRoom(roomID, data)

				// 缓存消息
				cache.CacheMessage(context.Background(), roomID, messageData)

				// 更新未读消息计数
				updateUnreadCounts(roomID, c.ID, msg.ID)
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket写入错误: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendMessage(msg WSMessage) {
	data, _ := json.Marshal(msg)
	select {
	case c.Send <- data:
	default:
	}
}

func generateConnID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

func getUserRoomIDs(userID uint) []uint {
	var roomIDs []uint
	database.GetDB().Model(&models.RoomMember{}).
		Where("user_id = ?", userID).
		Pluck("room_id", &roomIDs)
	return roomIDs
}

func isValidRoomMember(userID, roomID uint) bool {
	var count int64
	database.GetDB().Model(&models.RoomMember{}).
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Count(&count)
	return count > 0
}

func updateUnreadCounts(roomID, senderID, messageID uint) {
	var members []models.RoomMember
	database.GetDB().Where("room_id = ? AND user_id != ?", roomID, senderID).
		Find(&members)

	for _, member := range members {
		var unread models.UnreadMessage
		database.GetDB().Where("user_id = ? AND room_id = ?", member.UserID, roomID).
			FirstOrCreate(&unread)

		unread.Count++
		unread.LastMsgID = messageID
		database.GetDB().Save(&unread)
	}
}

func StartHub() {
	go hub.Run()
}
