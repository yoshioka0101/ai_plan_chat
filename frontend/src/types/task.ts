// Task type definitions

export type TaskStatus = 'todo' | 'in_progress' | 'done';

export interface Task {
  id: string;
  user_id: string;
  source: 'manual' | 'ai';
  title: string;
  description?: string;
  due_at?: string;
  status: TaskStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateTaskRequest {
  title: string;
  description?: string;
  due_at?: string;
  status?: TaskStatus;
}

export interface UpdateTaskRequest {
  title: string;
  description?: string;
  due_at?: string;
  status: TaskStatus;
}

export interface EditTaskRequest {
  title?: string;
  description?: string;
  due_at?: string;
  status?: TaskStatus;
}
