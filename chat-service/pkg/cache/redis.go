package cache

import (
	"chat-service/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	fmt.Println("Redis连接成功")
	return nil
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func Delete(ctx context.Context, keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

// SetUserOnline 设置用户在线状态
func SetUserOnline(ctx context.Context, userID uint, connID string, roomIDs []uint) error {
	key := fmt.Sprintf("user:online:%d", userID)
	data := map[string]interface{}{
		"conn_id":   connID,
		"room_ids":  roomIDs,
		"last_seen": time.Now().Unix(),
	}
	return Set(ctx, key, data, time.Hour*24)
}

// GetUserOnline 获取用户在线状态
func GetUserOnline(ctx context.Context, userID uint) (map[string]interface{}, error) {
	key := fmt.Sprintf("user:online:%d", userID)
	var data map[string]interface{}
	err := Get(ctx, key, &data)
	return data, err
}

// SetUserOffline 设置用户离线
func SetUserOffline(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:online:%d", userID)
	return Delete(ctx, key)
}

// GetRoomUsers 获取房间在线用户列表
func GetRoomUsers(ctx context.Context, roomID uint) ([]uint, error) {
	key := fmt.Sprintf("room:users:%d", roomID)
	userIDs, err := RedisClient.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result []uint
	for _, idStr := range userIDs {
		var id uint
		if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
			result = append(result, id)
		}
	}
	return result, nil
}

// AddUserToRoom 添加用户到房间
func AddUserToRoom(ctx context.Context, roomID, userID uint) error {
	key := fmt.Sprintf("room:users:%d", roomID)
	return RedisClient.SAdd(ctx, key, userID).Err()
}

// RemoveUserFromRoom 从房间移除用户
func RemoveUserFromRoom(ctx context.Context, roomID, userID uint) error {
	key := fmt.Sprintf("room:users:%d", roomID)
	return RedisClient.SRem(ctx, key, userID).Err()
}

// CacheMessage 缓存最近消息
func CacheMessage(ctx context.Context, roomID uint, message interface{}) error {
	key := fmt.Sprintf("room:messages:%d", roomID)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 使用列表存储最近100条消息
	pipe := RedisClient.Pipeline()
	pipe.LPush(ctx, key, data)
	pipe.LTrim(ctx, key, 0, 99)
	pipe.Expire(ctx, key, time.Hour*24)

	_, err = pipe.Exec(ctx)
	return err
}

// GetCachedMessages 获取缓存的消息
func GetCachedMessages(ctx context.Context, roomID uint, limit int) ([]interface{}, error) {
	key := fmt.Sprintf("room:messages:%d", roomID)
	data, err := RedisClient.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	var result []interface{}
	for _, item := range data {
		var msg interface{}
		if err := json.Unmarshal([]byte(item), &msg); err == nil {
			result = append(result, msg)
		}
	}
	return result, nil
}
