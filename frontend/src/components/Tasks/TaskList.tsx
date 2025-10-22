import { useState, useEffect } from 'react';
import { taskService } from '../../services/taskService';
import { TaskItem } from './TaskItem';
import { TaskForm } from './TaskForm';
import type { Task, CreateTaskRequest } from '../../types/task';

export const TaskList = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | undefined>();

  // Load tasks on component mount
  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await taskService.getTasks();
      setTasks(data);
    } catch (err) {
      setError('Failed to load tasks. Please try again.');
      console.error('Error loading tasks:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async (taskData: CreateTaskRequest) => {
    try {
      const newTask = await taskService.createTask(taskData);
      setTasks([...tasks, newTask]);
      setShowForm(false);
    } catch (err) {
      setError('Failed to create task. Please try again.');
      console.error('Error creating task:', err);
    }
  };

  const handleUpdateTask = async (taskData: CreateTaskRequest) => {
    if (!editingTask) return;

    try {
      // Use PUT endpoint for full update
      const updateData = {
        title: taskData.title,
        description: taskData.description,
        due_at: taskData.due_at,
        status: taskData.status || 'todo',
      };
      const updatedTask = await taskService.updateTask(editingTask.id, updateData);
      setTasks(tasks.map((t) => (t.id === updatedTask.id ? updatedTask : t)));
      setEditingTask(undefined);
      setShowForm(false);
    } catch (err) {
      setError('Failed to update task. Please try again.');
      console.error('Error updating task:', err);
    }
  };

  const handleDeleteTask = async (id: string) => {
    if (!confirm('Are you sure you want to delete this task?')) return;

    try {
      await taskService.deleteTask(id);
      setTasks(tasks.filter((t) => t.id !== id));
    } catch (err) {
      setError('Failed to delete task. Please try again.');
      console.error('Error deleting task:', err);
    }
  };

  const handleStatusChange = async (id: string, status: Task['status']) => {
    try {
      // Use PATCH endpoint for partial update
      const updatedTask = await taskService.editTask(id, { status });
      setTasks(tasks.map((t) => (t.id === updatedTask.id ? updatedTask : t)));
    } catch (err) {
      setError('Failed to update task status. Please try again.');
      console.error('Error updating status:', err);
    }
  };

  const handleEdit = (task: Task) => {
    setEditingTask(task);
    setShowForm(true);
  };

  const handleCancelForm = () => {
    setShowForm(false);
    setEditingTask(undefined);
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '40px', color: '#6b7280' }}>Loading tasks...</div>
    );
  }

  return (
    <div style={{ maxWidth: '900px', margin: '0 auto', padding: '24px' }}>
      <div style={{ marginBottom: '32px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
          <h1 style={{ margin: 0, fontSize: '28px', fontWeight: '700' }}>Tasks</h1>
          {!showForm && (
            <button
              onClick={() => setShowForm(true)}
              style={{
                padding: '10px 20px',
                borderRadius: '6px',
                border: 'none',
                backgroundColor: '#3b82f6',
                color: '#ffffff',
                cursor: 'pointer',
                fontSize: '14px',
                fontWeight: '500',
              }}
            >
              + New Task
            </button>
          )}
        </div>

        {error && (
          <div
            style={{
              padding: '12px',
              backgroundColor: '#fef2f2',
              border: '1px solid #fecaca',
              borderRadius: '6px',
              color: '#dc2626',
              marginBottom: '16px',
            }}
          >
            {error}
          </div>
        )}

        {showForm && (
          <TaskForm
            task={editingTask}
            onSubmit={editingTask ? handleUpdateTask : handleCreateTask}
            onCancel={handleCancelForm}
          />
        )}
      </div>

      {tasks.length === 0 ? (
        <div
          style={{
            textAlign: 'center',
            padding: '60px 20px',
            backgroundColor: '#f9fafb',
            borderRadius: '8px',
            color: '#6b7280',
          }}
        >
          <p style={{ margin: 0, fontSize: '16px' }}>No tasks yet. Create your first task to get started!</p>
        </div>
      ) : (
        <div>
          {tasks.map((task) => (
            <TaskItem
              key={task.id}
              task={task}
              onEdit={handleEdit}
              onDelete={handleDeleteTask}
              onStatusChange={handleStatusChange}
            />
          ))}
        </div>
      )}
    </div>
  );
};
