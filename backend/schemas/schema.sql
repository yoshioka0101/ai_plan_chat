-- AI Chat Task Interpretation System
-- Feature: 003-ai-chat-interpretation
-- Database: MySQL 8.0+

-- ============================================================================
-- users テーブル
-- ============================================================================
CREATE TABLE `users` (
  `id` char(36) NOT NULL COMMENT 'ユーザーID (UUID)',
  `nickname` varchar(255) NOT NULL COMMENT 'ニックネーム/表示名',
  `avatar` text NULL COMMENT 'アバター画像URL',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザープロフィール';

-- ============================================================================
-- user_auths テーブル
-- ============================================================================
CREATE TABLE `user_auths` (
  `id` char(36) NOT NULL COMMENT '認証ID (UUID)',
  `user_id` char(36) NOT NULL COMMENT 'ユーザーID',
  `identity_type` varchar(50) NOT NULL COMMENT '認証タイプ（google, email, github等）',
  `identifier` varchar(255) NOT NULL COMMENT '識別子（メールアドレス、ユーザー名等）',
  `credential` text NULL COMMENT '認証情報（トークン、認証コード等）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_auths_identity` (`user_id`, `identity_type`, `identifier`),
  KEY `idx_user_auths_identifier` (`identity_type`, `identifier`),
  CONSTRAINT `fk_user_auths_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザー認証情報';

-- ============================================================================
-- ai_interpretations テーブル
-- ============================================================================
CREATE TABLE `ai_interpretations` (
  `id` char(36) NOT NULL COMMENT 'AI解釈ID (UUID)',
  `user_id` char(36) NOT NULL COMMENT 'ユーザーID',
  `input_text` text NOT NULL COMMENT 'ユーザーが入力した自然言語テキスト',
  `structured_result` json NOT NULL COMMENT 'AI解析結果のJSON構造',
  `ai_model` varchar(100) NOT NULL DEFAULT 'gemini-flash' COMMENT '使用AIモデル名',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '解析実行日時',
  PRIMARY KEY (`id`),
  KEY `idx_ai_interpretations_user_created` (`user_id`, `created_at` DESC),
  KEY `idx_ai_interpretations_type` ((cast(json_unquote(json_extract(`structured_result`,_utf8mb4'$.type')) as char(50) charset utf8mb4))),
  CONSTRAINT `fk_ai_interpretations_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI解析履歴';

-- ============================================================================
-- tasks テーブル
-- ============================================================================
CREATE TABLE `tasks` (
  `id` char(36) NOT NULL COMMENT 'タスクID (UUID)',
  `user_id` char(36) NOT NULL COMMENT 'ユーザーID',
  `title` varchar(255) NOT NULL COMMENT 'タスクのタイトル',
  `description` text NULL COMMENT '詳細説明',
  `due_at` timestamp NULL COMMENT '期限日時',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT 'ステータス',
  `source` varchar(20) NOT NULL DEFAULT 'manual' COMMENT '作成元',
  `ai_interpretation_id` char(36) NULL COMMENT '元のAI解釈ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最終更新日時',
  PRIMARY KEY (`id`),
  KEY `idx_tasks_user_due` (`user_id`, `due_at`),
  KEY `idx_tasks_user_status` (`user_id`, `status`),
  KEY `idx_tasks_user_created` (`user_id`, `created_at` DESC),
  KEY `fk_tasks_ai_interpretation` (`ai_interpretation_id`),
  CONSTRAINT `fk_tasks_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tasks_ai_interpretation` FOREIGN KEY (`ai_interpretation_id`) REFERENCES `ai_interpretations` (`id`) ON DELETE SET NULL,
  CONSTRAINT `chk_tasks_status` CHECK (`status` IN ('pending', 'completed', 'cancelled')),
  CONSTRAINT `chk_tasks_source` CHECK (`source` IN ('ai', 'manual'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='タスク';

