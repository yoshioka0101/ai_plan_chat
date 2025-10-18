-- Migration: Create tasks table
-- Created: 2025-10-16
-- Description: Task API用のtasksテーブルを作成

CREATE TABLE IF NOT EXISTS tasks (
    id CHAR(36) PRIMARY KEY COMMENT 'タスクID (UUID)',
    user_id CHAR(36) NOT NULL COMMENT 'ユーザーID',
    title TEXT NOT NULL COMMENT 'タスクタイトル',
    description TEXT COMMENT 'タスク詳細',
    due_at TIMESTAMP NULL COMMENT '期限日時',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT 'ステータス: pending/in_progress/completed',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_due_at (due_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='タスク管理テーブル';
