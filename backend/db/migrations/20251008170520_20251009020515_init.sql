f-- Create "users" table
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `nickname` text NOT NULL,
  `avatar` text NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "events" table
CREATE TABLE `events` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `title` varchar(500) NOT NULL,
  `description` text NULL,
  `starts_at` timestamp NOT NULL,
  `ends_at` timestamp NOT NULL,
  `location` varchar(500) NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_events_starts_at` (`starts_at`),
  INDEX `idx_events_user_id` (`user_id`),
  CONSTRAINT `events_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "interpretations" table
CREATE TABLE `interpretations` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `name` varchar(255) NULL,
  `description` text NULL,
  `context` text NULL,
  `input_text` text NOT NULL,
  `structured_result` json NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_interpretations_user_id` (`user_id`),
  CONSTRAINT `interpretations_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tags" table
CREATE TABLE `tags` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `name` varchar(100) NOT NULL,
  `color` varchar(20) NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_tags_user_id` (`user_id`),
  UNIQUE INDEX `tags_user_id_name_key` (`user_id`, `name`),
  CONSTRAINT `tags_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tag_events" table
CREATE TABLE `tag_events` (
  `id` varchar(36) NOT NULL,
  `tag_id` varchar(36) NOT NULL,
  `event_id` varchar(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_tag_events_event_id` (`event_id`),
  INDEX `idx_tag_events_tag_id` (`tag_id`),
  UNIQUE INDEX `tag_events_tag_id_event_id_key` (`tag_id`, `event_id`),
  CONSTRAINT `tag_events_event_id_fkey` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `tag_events_tag_id_fkey` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tasks" table
CREATE TABLE `tasks` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `title` varchar(500) NOT NULL,
  `description` text NULL,
  `due_at` timestamp NULL,
  `status` varchar(20) NOT NULL DEFAULT "pending",
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_tasks_due_at` (`due_at`),
  INDEX `idx_tasks_status` (`status`),
  INDEX `idx_tasks_user_id` (`user_id`),
  CONSTRAINT `tasks_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tag_tasks" table
CREATE TABLE `tag_tasks` (
  `id` varchar(36) NOT NULL,
  `tag_id` varchar(36) NOT NULL,
  `task_id` varchar(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_tag_tasks_tag_id` (`tag_id`),
  INDEX `idx_tag_tasks_task_id` (`task_id`),
  UNIQUE INDEX `tag_tasks_tag_id_task_id_key` (`tag_id`, `task_id`),
  CONSTRAINT `tag_tasks_tag_id_fkey` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `tag_tasks_task_id_fkey` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_auths" table
CREATE TABLE `user_auths` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `identity_type` varchar(50) NOT NULL,
  `identifier` varchar(255) NOT NULL,
  `credential` text NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_auths_user_id` (`user_id`),
  UNIQUE INDEX `user_auths_user_id_identity_type_identifier_key` (`user_id`, `identity_type`, `identifier`),
  CONSTRAINT `user_auths_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
