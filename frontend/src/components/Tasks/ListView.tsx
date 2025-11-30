import { useState } from 'react';
import { format } from 'date-fns';
import type { Task, TaskStatus } from '../../types/task';

interface ListViewProps {
  tasks: Task[];
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: TaskStatus) => void;
  onNewTaskClick: () => void;
  onUpdate: (id: string, updates: { title?: string; description?: string; due_at?: string }) => void;
}

export const ListView = ({ tasks, onDelete, onStatusChange, onNewTaskClick, onUpdate }: ListViewProps) => {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editTitle, setEditTitle] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editDueDate, setEditDueDate] = useState('');
  const startEdit = (task: Task) => {
    setEditingId(task.id);
    setEditTitle(task.title);
    setEditDescription(task.description || '');
    setEditDueDate(task.due_at ? task.due_at.split('T')[0] : '');
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditTitle('');
    setEditDescription('');
    setEditDueDate('');
  };

  const saveEdit = (taskId: string) => {
    onUpdate(taskId, {
      title: editTitle,
      description: editDescription || undefined,
      due_at: editDueDate || undefined,
    });
    cancelEdit();
  };

  const getStatusColor = (status: TaskStatus) => {
    switch (status) {
      case 'todo':
        return { bg: '#fef3c7', text: '#92400e' }; // amber - darker text
      case 'in_progress':
        return { bg: '#dbeafe', text: '#1e40af' }; // blue - darker text
      case 'done':
        return { bg: '#d1fae5', text: '#065f46' }; // green - darker text
      default:
        return { bg: '#f3f4f6', text: '#374151' }; // gray - darker text
    }
  };

  return (
    <div style={{ width: '100%' }}>
      {/* 新規作成ボタン */}
      <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '16px' }}>
        <button
          onClick={onNewTaskClick}
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

      {/* ヘッダー */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: '3fr 2fr 1fr 1fr 180px',
          gap: '16px',
          padding: '16px 20px',
          backgroundColor: '#f9fafb',
          borderRadius: '8px',
          marginBottom: '12px',
          fontWeight: '600',
          fontSize: '14px',
          color: '#6b7280',
        }}
      >
        <div>タイトル</div>
        <div>説明</div>
        <div>期限</div>
        <div>ステータス</div>
        <div style={{ textAlign: 'right' }}>操作</div>
      </div>

      {/* タスクリスト */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
        {tasks.length === 0 ? (
          <div
            style={{
              textAlign: 'center',
              padding: '60px 20px',
              color: '#9ca3af',
              fontSize: '14px',
            }}
          >
            タスクがありません
          </div>
        ) : (
          tasks.map((task) => {
            const statusColor = getStatusColor(task.status);
            const isEditing = editingId === task.id;

            return (
              <div
                key={task.id}
                style={{
                  display: 'grid',
                  gridTemplateColumns: '3fr 2fr 1fr 1fr 180px',
                  gap: '16px',
                  padding: '16px 20px',
                  backgroundColor: isEditing ? '#fef3c7' : '#ffffff',
                  border: isEditing ? '2px solid #f59e0b' : '1px solid #e5e7eb',
                  borderRadius: '8px',
                  alignItems: 'center',
                  transition: 'all 0.2s',
                }}
                onMouseEnter={(e) => {
                  if (!isEditing) {
                    e.currentTarget.style.boxShadow = '0 2px 4px rgba(0, 0, 0, 0.1)';
                  }
                }}
                onMouseLeave={(e) => {
                  if (!isEditing) {
                    e.currentTarget.style.boxShadow = 'none';
                  }
                }}
              >
                {/* タイトル */}
                <div>
                  {isEditing ? (
                    <input
                      type="text"
                      value={editTitle}
                      onChange={(e) => setEditTitle(e.target.value)}
                      style={{
                        width: '100%',
                        padding: '8px',
                        fontSize: '14px',
                        fontWeight: '600',
                        border: '1px solid #d1d5db',
                        borderRadius: '6px',
                        backgroundColor: '#ffffff',
                        boxSizing: 'border-box',
                      }}
                      placeholder="タイトル"
                      autoFocus
                    />
                  ) : (
                    <div
                      style={{
                        fontWeight: '600',
                        fontSize: '14px',
                        color: '#1f2937',
                        cursor: 'pointer',
                      }}
                      onClick={() => startEdit(task)}
                    >
                      {task.title}
                    </div>
                  )}
                </div>

                {/* 説明 */}
                <div>
                  {isEditing ? (
                    <input
                      type="text"
                      value={editDescription}
                      onChange={(e) => setEditDescription(e.target.value)}
                      style={{
                        width: '100%',
                        padding: '8px',
                        fontSize: '13px',
                        border: '1px solid #d1d5db',
                        borderRadius: '6px',
                        backgroundColor: '#ffffff',
                        boxSizing: 'border-box',
                      }}
                      placeholder="説明"
                    />
                  ) : (
                    <div
                      style={{
                        fontSize: '13px',
                        color: '#6b7280',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        cursor: 'pointer',
                      }}
                      onClick={() => startEdit(task)}
                    >
                      {task.description || '-'}
                    </div>
                  )}
                </div>

                {/* 期限 */}
                <div>
                  {isEditing ? (
                    <input
                      type="date"
                      value={editDueDate}
                      onChange={(e) => setEditDueDate(e.target.value)}
                      style={{
                        width: '100%',
                        padding: '6px',
                        fontSize: '13px',
                        border: '1px solid #d1d5db',
                        borderRadius: '6px',
                        backgroundColor: '#ffffff',
                        boxSizing: 'border-box',
                      }}
                    />
                  ) : (
                    <div style={{ fontSize: '13px', color: '#6b7280' }}>
                      {task.due_at ? format(new Date(task.due_at), 'yyyy/MM/dd') : '-'}
                    </div>
                  )}
                </div>

                {/* ステータス */}
                <div>
                  <select
                    value={task.status}
                    onChange={(e) => onStatusChange(task.id, e.target.value as TaskStatus)}
                    style={{
                      padding: '6px 10px',
                      fontSize: '12px',
                      fontWeight: '600',
                      border: '1px solid #d1d5db',
                      borderRadius: '6px',
                      backgroundColor: statusColor.bg,
                      color: statusColor.text,
                      cursor: 'pointer',
                      width: '100%',
                    }}
                  >
                    <option value="todo">To Do</option>
                    <option value="in_progress">In Progress</option>
                    <option value="done">Done</option>
                  </select>
                </div>

                {/* 操作ボタン */}
                <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                  {isEditing ? (
                    <>
                      <button
                        onClick={() => saveEdit(task.id)}
                        style={{
                          padding: '6px 12px',
                          fontSize: '12px',
                          fontWeight: '600',
                          border: 'none',
                          borderRadius: '6px',
                          backgroundColor: '#10b981',
                          color: '#ffffff',
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = '#059669';
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = '#10b981';
                        }}
                      >
                        ✓ 保存
                      </button>
                      <button
                        onClick={cancelEdit}
                        style={{
                          padding: '6px 12px',
                          fontSize: '12px',
                          fontWeight: '500',
                          border: '1px solid #d1d5db',
                          borderRadius: '6px',
                          backgroundColor: '#ffffff',
                          color: '#374151',
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = '#f9fafb';
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = '#ffffff';
                        }}
                      >
                        ✕ キャンセル
                      </button>
                    </>
                  ) : (
                    <>
                      <button
                        onClick={() => startEdit(task)}
                        style={{
                          padding: '6px 12px',
                          fontSize: '12px',
                          fontWeight: '500',
                          border: '1px solid #d1d5db',
                          borderRadius: '6px',
                          backgroundColor: '#ffffff',
                          color: '#374151',
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = '#f9fafb';
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = '#ffffff';
                        }}
                      >
                        ✏️ 編集
                      </button>
                      <button
                        onClick={() => onDelete(task.id)}
                        style={{
                          padding: '6px 12px',
                          fontSize: '12px',
                          fontWeight: '600',
                          border: '1px solid #dc2626',
                          borderRadius: '6px',
                          backgroundColor: '#dc2626',
                          color: '#ffffff',
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = '#b91c1c';
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = '#dc2626';
                        }}
                      >
                        削除
                      </button>
                    </>
                  )}
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
};
