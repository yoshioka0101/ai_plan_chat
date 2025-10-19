-- Create "tasks" table
CREATE TABLE `tasks` (
  `id` char(36) NOT NULL COMMENT "タスクID (UUID)",
  `title` varchar(500) NOT NULL COMMENT "タスクタイトル",
  `description` text NULL COMMENT "タスク詳細",
  `due_at` timestamp NULL COMMENT "期限日時",
  `status` varchar(20) NOT NULL DEFAULT "pending" COMMENT "ステータス（pending/in_progress/completed）",
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`id`),
  INDEX `idx_tasks_created_at` (`created_at`),
  INDEX `idx_tasks_due_at` (`due_at`),
  INDEX `idx_tasks_status` (`status`),
  CONSTRAINT `tasks_chk_1` CHECK (`status` in (_utf8mb4'pending',_utf8mb4'in_progress',_utf8mb4'completed'))
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "タスク";
