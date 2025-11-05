package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// SanitizeInput 清理输入字符串
func SanitizeInput(input string) string {
	// 去除首尾空格
	input = strings.TrimSpace(input)

	// 转义HTML特殊字符
	input = strings.ReplaceAll(input, "&", "&amp;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")

	return input
}

// FormatMessage 格式化消息内容
func FormatMessage(content string) string {
	// 简单的消息格式化
	// 可以根据需要扩展：支持表情、链接解析等
	content = strings.TrimSpace(content)

	// 限制消息长度
	if len(content) > 1000 {
		content = content[:1000] + "..."
	}

	return content
}

// GetRoomType 获取房间类型描述
func GetRoomType(roomType string) string {
	switch roomType {
	case "single":
		return "单人聊天"
	case "group":
		return "群组聊天"
	default:
		return "未知类型"
	}
}

// ValidateRoomType 验证房间类型
func ValidateRoomType(roomType string) bool {
	validTypes := []string{"single", "group"}
	for _, t := range validTypes {
		if t == roomType {
			return true
		}
	}
	return false
}

// FormatResponse 统一格式化API响应
func FormatResponse(code int, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// CalculatePagination 计算分页信息
func CalculatePagination(page, pageSize int, total int64) PaginationInfo {
	totalPage := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginationInfo{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}

// Contains 检查切片是否包含元素
func Contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveFromSlice 从切片中移除元素
func RemoveFromSlice[T comparable](slice []T, item T) []T {
	var result []T
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// TimeAgo 计算时间差显示
func TimeAgo(t interface{}) string {
	// 这里可以实现时间差计算逻辑
	// 例如：1分钟前、2小时前等
	return fmt.Sprintf("%v", t)
}
