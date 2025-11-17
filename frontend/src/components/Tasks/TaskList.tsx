import { useState, useEffect, useCallback } from 'react';
import { taskService } from '../../services/taskService';
import { useAuth } from '../../hooks/useAuth';
import { KanbanColumn } from './KanbanColumn';
import { ListView } from './ListView';
import { CalendarView } from './CalendarView';
import { AIView } from './AIView';
import { Sidebar } from '../Sidebar/Sidebar';
import { Modal } from '../Modal/Modal';
import { TaskForm } from './TaskForm';
import type { Task, CreateTaskRequest } from '../../types/task';

type ViewMode = 'kanban' | 'list' | 'calendar' | 'ai';

export const TaskList = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | undefined>();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [viewMode, setViewMode] = useState<ViewMode>('kanban');

  const { token, isAuthenticated } = useAuth();

  const loadTasks = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await taskService.getTasks();
      setTasks(Array.isArray(data) ? data : []);
    } catch (err) {
      setError('Failed to load tasks. Please try again.');
      console.error('Error loading tasks:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load tasks on component mount and when authentication state changes
  useEffect(() => {
    if (isAuthenticated && token) {
      loadTasks();
    }
  }, [isAuthenticated, token, loadTasks]);

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
    if (!confirm('削除してもよろしいでしょうか')) return;

    try {
      await taskService.deleteTask(id);
      setTasks(tasks.filter((t) => t.id !== id));

      // 削除したタスクを編集中だった場合、編集状態をクリア
      if (editingTask && editingTask.id === id) {
        setEditingTask(undefined);
        setShowForm(false);
      }
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

  const handleInlineUpdate = async (id: string, updates: { title?: string; description?: string; due_at?: string }) => {
    try {
      // Use PATCH endpoint for partial update
      const updatedTask = await taskService.editTask(id, updates);
      setTasks(tasks.map((t) => (t.id === updatedTask.id ? updatedTask : t)));
    } catch (err) {
      setError('Failed to update task. Please try again.');
      console.error('Error updating task:', err);
    }
  };

  const handleEdit = async (task: Task) => {
    try {
      // タスク詳細をAPIから取得
      const taskDetail = await taskService.getTask(task.id);
      setEditingTask(taskDetail);
      setShowForm(true);
    } catch (err) {
      setError('Failed to load task details. Please try again.');
      console.error('Error loading task details:', err);
    }
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

  // タスクをステータスごとにグループ化
  const todoTasks = tasks.filter((task) => task.status === 'todo');
  const inProgressTasks = tasks.filter((task) => task.status === 'in_progress');
  const doneTasks = tasks.filter((task) => task.status === 'done');

  return (
    <>
      {/* サイドバー */}
      <Sidebar
        onNewTaskClick={() => setShowForm(true)}
        isOpen={sidebarOpen}
        onToggle={setSidebarOpen}
        viewMode={viewMode}
        onViewChange={setViewMode}
      />

      {/* メインコンテンツ */}
      <div
        style={{
          minHeight: '100vh',
          marginLeft: sidebarOpen ? '320px' : '0',
          padding: '32px',
          backgroundColor: '#f3f4f6',
          transition: 'margin-left 0.3s ease',
        }}
      >
        <div style={{ maxWidth: '1600px', margin: '0 auto', width: '100%' }}>
          {/* ヘッダー */}
          <div style={{ marginBottom: '32px' }}>
            <div>
              <h1 style={{ margin: '0 0 8px 0', fontSize: '32px', fontWeight: '700', color: '#1f2937' }}>
                Task Board
              </h1>
              <p style={{ margin: 0, fontSize: '14px', color: '#6b7280' }}>
                Manage your tasks with Kanban style
              </p>
            </div>

            {error && (
              <div
                style={{
                  padding: '12px 16px',
                  backgroundColor: '#fef2f2',
                  border: '1px solid #fecaca',
                  borderRadius: '8px',
                  color: '#dc2626',
                  marginTop: '16px',
                  fontSize: '14px',
                }}
              >
                {error}
              </div>
            )}
          </div>

        {/* ビューコンテンツ */}
        {tasks.length === 0 && !showForm ? (
          <div
            style={{
              textAlign: 'center',
              padding: '80px 20px',
              backgroundColor: '#ffffff',
              borderRadius: '12px',
              border: '2px dashed #e5e7eb',
            }}
          >
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>📋</div>
            <p style={{ margin: '0 0 8px 0', fontSize: '18px', fontWeight: '600', color: '#1f2937' }}>
              No tasks yet
            </p>
            <p style={{ margin: 0, fontSize: '14px', color: '#6b7280' }}>
              Create your first task to get started!
            </p>
          </div>
        ) : viewMode === 'kanban' ? (
          <>
            {/* Kanban 新規作成ボタン */}
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '16px' }}>
              <button
                onClick={() => setShowForm(true)}
                style={{
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: '600',
                  border: 'none',
                  borderRadius: '8px',
                  backgroundColor: '#3b82f6',
                  color: '#ffffff',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
                  transition: 'all 0.2s',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = '#2563eb';
                  e.currentTarget.style.transform = 'translateY(-1px)';
                  e.currentTarget.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.15)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = '#3b82f6';
                  e.currentTarget.style.transform = 'translateY(0)';
                  e.currentTarget.style.boxShadow = '0 2px 4px rgba(0, 0, 0, 0.1)';
                }}
              >
                <span style={{ fontSize: '18px' }}>+</span>
                新規タスク作成
              </button>
            </div>
            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(3, 1fr)',
                gap: '24px',
              }}
            >
            <KanbanColumn
              title="To Do"
              status="todo"
              tasks={todoTasks}
              onEdit={handleEdit}
              onDelete={handleDeleteTask}
              onStatusChange={handleStatusChange}
            />
            <KanbanColumn
              title="In Progress"
              status="in_progress"
              tasks={inProgressTasks}
              onEdit={handleEdit}
              onDelete={handleDeleteTask}
              onStatusChange={handleStatusChange}
            />
            <KanbanColumn
              title="Done"
              status="done"
              tasks={doneTasks}
              onEdit={handleEdit}
              onDelete={handleDeleteTask}
              onStatusChange={handleStatusChange}
            />
          </div>
          </>
        ) : viewMode === 'list' ? (
          <ListView
            tasks={tasks}
            onDelete={handleDeleteTask}
            onStatusChange={handleStatusChange}
            onNewTaskClick={() => setShowForm(true)}
            onUpdate={handleInlineUpdate}
          />
        ) : viewMode === 'calendar' ? (
          <CalendarView
            tasks={tasks}
            onEdit={handleEdit}
            onDelete={handleDeleteTask}
            onStatusChange={handleStatusChange}
            onNewTaskClick={() => setShowForm(true)}
          />
        ) : (
          <AIView
            onNewTaskClick={() => setShowForm(true)}
          />
        )}
        </div>
      </div>

      {/* タスク編集・作成モーダル */}
      <Modal
        isOpen={showForm}
        onClose={handleCancelForm}
        title={editingTask ? 'タスクを編集' : '新規タスク作成'}
      >
        <TaskForm
          task={editingTask}
          onSubmit={editingTask ? handleUpdateTask : handleCreateTask}
          onCancel={handleCancelForm}
        />
      </Modal>
    </>
  );
};
