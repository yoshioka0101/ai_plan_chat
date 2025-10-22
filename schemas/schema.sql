-- Task Management Schema
-- タスク管理に必要な最小限のスキーマ

-- tasks（タスク）
CREATE TABLE IF NOT EXISTS tasks (
  id CHAR(36) PRIMARY KEY COMMENT 'タスクID (UUID)',
  title VARCHAR(500) NOT NULL COMMENT 'タスクタイトル',
  description TEXT COMMENT 'タスク詳細',
  due_at TIMESTAMP NULL COMMENT '期限日時',
  status VARCHAR(20) NOT NULL DEFAULT 'todo' COMMENT 'ステータス（todo/in_progress/done）',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

  INDEX idx_tasks_status (status),
  INDEX idx_tasks_due_at (due_at),
  INDEX idx_tasks_created_at (created_at),
  CHECK (status IN ('todo', 'in_progress', 'done'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='タスク';
