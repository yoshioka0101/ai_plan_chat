import { apiClient } from './api';
import type { Task, CreateTaskRequest, UpdateTaskRequest } from '../types/task';

export const taskService = {
  // Get all tasks
  async getTasks(): Promise<Task[]> {
    const response = await apiClient.get<Task[]>('/tasks');
    return response.data;
  },

  // Get a single task by ID
  async getTask(id: string): Promise<Task> {
    const response = await apiClient.get<Task>(`/tasks/${id}`);
    return response.data;
  },

  // Create a new task
  async createTask(task: CreateTaskRequest): Promise<Task> {
    const response = await apiClient.post<Task>('/tasks', task);
    return response.data;
  },

  // Update an existing task
  async updateTask(id: string, task: UpdateTaskRequest): Promise<Task> {
    const response = await apiClient.put<Task>(`/tasks/${id}`, task);
    return response.data;
  },

  // Delete a task
  async deleteTask(id: string): Promise<void> {
    await apiClient.delete(`/tasks/${id}`);
  },
};
