-- Update task status values to match OpenAPI definition
-- Change from pending/in_progress/completed to todo/in_progress/done

-- Update existing records
UPDATE tasks SET status = 'todo' WHERE status = 'pending';
UPDATE tasks SET status = 'done' WHERE status = 'completed';

-- Update table constraint
ALTER TABLE tasks DROP CONSTRAINT tasks_chk_1;
ALTER TABLE tasks ADD CONSTRAINT tasks_chk_1 CHECK (status IN ('todo', 'in_progress', 'done'));

-- Update default value
ALTER TABLE tasks ALTER COLUMN status SET DEFAULT 'todo';
