-- Modify "ai_interpretations" table
ALTER TABLE `ai_interpretations` ADD COLUMN `original_result` json NULL COMMENT "AI提案の原本（レビュー前）" AFTER `structured_result`;
-- Create "interpretation_items" table
CREATE TABLE `interpretation_items` (
  `id` char(36) NOT NULL COMMENT "アイテムID (UUID)",
  `interpretation_id` char(36) NOT NULL COMMENT "AI解釈ID",
  `item_index` int NOT NULL COMMENT "結果内のindex",
  `resource_type` varchar(20) NOT NULL COMMENT "リソースタイプ (task/event/wallet)",
  `resource_id` char(36) NULL COMMENT "作成済みリソースID",
  `status` varchar(20) NOT NULL DEFAULT "pending" COMMENT "ステータス (pending/created)",
  `data` json NOT NULL COMMENT "編集後のアイテム内容",
  `original_data` json NOT NULL COMMENT "AI提案の原本",
  `reviewed_at` timestamp NULL COMMENT "レビュー日時",
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`id`),
  INDEX `idx_interpretation_items_interpretation` (`interpretation_id`, `item_index`),
  INDEX `idx_interpretation_items_resource` (`resource_type`, `resource_id`),
  INDEX `idx_interpretation_items_status` (`interpretation_id`, `status`),
  CONSTRAINT `fk_interpretation_items_interpretation` FOREIGN KEY (`interpretation_id`) REFERENCES `ai_interpretations` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `chk_interpretation_items_resource_type` CHECK (`resource_type` in (_utf8mb4'task',_utf8mb4'event',_utf8mb4'wallet')),
  CONSTRAINT `chk_interpretation_items_status` CHECK (`status` in (_utf8mb4'pending',_utf8mb4'created'))
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "AI解釈アイテム（レビュー対象）";
