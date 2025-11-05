package database

import (
	"chat-service/internal/config"
	"chat-service/internal/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.User{},
		&models.ChatRoom{},
		&models.RoomMember{},
		&models.Message{},
		&models.UnreadMessage{},
		&models.OnlineUser{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库连接成功并完成迁移")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
