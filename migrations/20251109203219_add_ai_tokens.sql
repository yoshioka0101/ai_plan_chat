-- Modify "ai_interpretations" table
ALTER TABLE `ai_interpretations` ADD COLUMN `ai_prompt_tokens` int NULL COMMENT "入力トークン数" AFTER `ai_model`, ADD COLUMN `ai_completion_tokens` int NULL COMMENT "出力トークン数" AFTER `ai_prompt_tokens`;
-- Modify "tasks" table
ALTER TABLE `tasks` DROP CONSTRAINT `tasks_chk_1`, ADD CONSTRAINT `chk_tasks_status` CHECK (`status` in (_utf8mb4'todo',_utf8mb4'in_progress',_utf8mb4'done')), MODIFY COLUMN `status` varchar(20) NOT NULL DEFAULT "todo" COMMENT "ステータス（todo/in_progress/done）";
