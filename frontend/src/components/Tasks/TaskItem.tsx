import { format } from 'date-fns';
import type { Task } from '../../types/task';

interface TaskItemProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: Task['status']) => void;
}

export const TaskItem = ({ task, onEdit, onDelete, onStatusChange }: TaskItemProps) => {
  const getStatusColor = (status: Task['status']) => {
    switch (status) {
      case 'done':
        return '#4ade80';
      case 'in_progress':
        return '#60a5fa';
      case 'todo':
        return '#fbbf24';
      default:
        return '#94a3b8';
    }
  };

  return (
    <div
      style={{
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        padding: '16px',
        marginBottom: '12px',
        backgroundColor: '#ffffff',
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
        <div style={{ flex: 1 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '8px' }}>
            <h3 style={{ margin: 0, fontSize: '18px', fontWeight: '600' }}>{task.title}</h3>
          </div>

          {task.description && (
            <p style={{ margin: '0 0 12px 0', color: '#6b7280', fontSize: '14px' }}>
              {task.description}
            </p>
          )}

          <div style={{ display: 'flex', gap: '16px', fontSize: '13px', color: '#9ca3af' }}>
            <div>
              Status:
              <select
                value={task.status}
                onChange={(e) => onStatusChange(task.id, e.target.value as Task['status'])}
                style={{
                  marginLeft: '8px',
                  padding: '4px 8px',
                  borderRadius: '4px',
                  border: '1px solid #d1d5db',
                  backgroundColor: getStatusColor(task.status) + '20',
                  color: getStatusColor(task.status),
                  fontWeight: '500',
                  fontSize: '12px',
                }}
              >
                <option value="todo">To Do</option>
                <option value="in_progress">In Progress</option>
                <option value="done">Done</option>
              </select>
            </div>

            {task.due_at && (
              <div>Due: {format(new Date(task.due_at), 'MMM d, yyyy')}</div>
            )}
          </div>
        </div>

        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={() => onEdit(task)}
            style={{
              padding: '6px 12px',
              borderRadius: '6px',
              border: '1px solid #d1d5db',
              backgroundColor: '#ffffff',
              color: '#374151',
              cursor: 'pointer',
              fontSize: '14px',
            }}
          >
            Edit
          </button>
          <button
            onClick={() => onDelete(task.id)}
            style={{
              padding: '6px 12px',
              borderRadius: '6px',
              border: '1px solid #fecaca',
              backgroundColor: '#fef2f2',
              color: '#dc2626',
              cursor: 'pointer',
              fontSize: '14px',
            }}
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
};
