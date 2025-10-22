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
      case 'completed':
        return '#4ade80';
      case 'in_progress':
        return '#60a5fa';
      case 'cancelled':
        return '#94a3b8';
      default:
        return '#fbbf24';
    }
  };

  const getPriorityColor = (priority?: Task['priority']) => {
    switch (priority) {
      case 'high':
        return '#ef4444';
      case 'medium':
        return '#f59e0b';
      case 'low':
        return '#10b981';
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
            {task.priority && (
              <span
                style={{
                  fontSize: '12px',
                  padding: '2px 8px',
                  borderRadius: '4px',
                  backgroundColor: getPriorityColor(task.priority) + '20',
                  color: getPriorityColor(task.priority),
                  fontWeight: '500',
                }}
              >
                {task.priority}
              </span>
            )}
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
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>

            {task.due_date && (
              <div>Due: {format(new Date(task.due_date), 'MMM d, yyyy')}</div>
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
