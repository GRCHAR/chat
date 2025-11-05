-- 聊天服务数据库初始化脚本

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS chat_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE chat_service;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    username varchar(50) NOT NULL,
    nickname varchar(100) NOT NULL,
    avatar varchar(255) DEFAULT NULL,
    email varchar(100) NOT NULL,
    password varchar(255) NOT NULL,
    status varchar(20) DEFAULT 'active',
    created_at datetime(3) NULL,
    updated_at datetime(3) NULL,
    deleted_at datetime(3) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_users_username (username),
    UNIQUE KEY idx_users_email (email),
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 聊天室表
CREATE TABLE IF NOT EXISTS chat_rooms (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    name varchar(100) NOT NULL,
    description varchar(500) DEFAULT NULL,
    type varchar(20) DEFAULT 'group',
    avatar varchar(255) DEFAULT NULL,
    owner_id bigint unsigned NOT NULL,
    max_members int DEFAULT '100',
    created_at datetime(3) NULL,
    updated_at datetime(3) NULL,
    deleted_at datetime(3) NULL,
    PRIMARY KEY (id),
    KEY idx_chat_rooms_owner_id (owner_id),
    KEY idx_chat_rooms_deleted_at (deleted_at),
    CONSTRAINT fk_chat_rooms_owner FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 聊天室成员表
CREATE TABLE IF NOT EXISTS room_members (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    room_id bigint unsigned NOT NULL,
    user_id bigint unsigned NOT NULL,
    role varchar(20) DEFAULT 'member',
    joined_at datetime(3) NULL,
    last_read datetime(3) NULL,
    created_at datetime(3) NULL,
    updated_at datetime(3) NULL,
    deleted_at datetime(3) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_room_members_room_user (room_id, user_id),
    KEY idx_room_members_user_id (user_id),
    KEY idx_room_members_deleted_at (deleted_at),
    CONSTRAINT fk_room_members_room FOREIGN KEY (room_id) REFERENCES chat_rooms (id) ON DELETE CASCADE,
    CONSTRAINT fk_room_members_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    room_id bigint unsigned NOT NULL,
    sender_id bigint unsigned NOT NULL,
    content longtext NOT NULL,
    type varchar(20) DEFAULT 'text',
    reply_to_id bigint unsigned DEFAULT NULL,
    is_deleted tinyint(1) DEFAULT '0',
    created_at datetime(3) NULL,
    updated_at datetime(3) NULL,
    deleted_at datetime(3) NULL,
    PRIMARY KEY (id),
    KEY idx_messages_room_id (room_id),
    KEY idx_messages_sender_id (sender_id),
    KEY idx_messages_reply_to_id (reply_to_id),
    KEY idx_messages_deleted_at (deleted_at),
    CONSTRAINT fk_messages_room FOREIGN KEY (room_id) REFERENCES chat_rooms (id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_reply_to FOREIGN KEY (reply_to_id) REFERENCES messages (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 未读消息表
CREATE TABLE IF NOT EXISTS unread_messages (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    user_id bigint unsigned NOT NULL,
    room_id bigint unsigned NOT NULL,
    count int DEFAULT '0',
    last_msg_id bigint unsigned NOT NULL,
    updated_at datetime(3) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_unread_messages_user_room (user_id, room_id),
    KEY idx_unread_messages_user_id (user_id),
    KEY idx_unread_messages_room_id (room_id),
    CONSTRAINT fk_unread_messages_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_unread_messages_room FOREIGN KEY (room_id) REFERENCES chat_rooms (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 在线用户表
CREATE TABLE IF NOT EXISTS online_users (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    user_id bigint unsigned NOT NULL,
    conn_id varchar(255) NOT NULL,
    room_ids longtext,
    last_seen datetime(3) NULL,
    created_at datetime(3) NULL,
    updated_at datetime(3) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_online_users_conn_id (conn_id),
    KEY idx_online_users_user_id (user_id),
    CONSTRAINT fk_online_users_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_messages_room_created ON messages (room_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_room_members_user_joined ON room_members (user_id, joined_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_rooms_type_owner ON chat_rooms (type, owner_id);

-- 插入示例数据（可选）
-- INSERT INTO users (username, nickname, email, password, status, created_at, updated_at) VALUES
-- ('admin', '系统管理员', 'admin@chat.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
