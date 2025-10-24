import { format } from 'date-fns';
import type { Task, TaskStatus } from '../../types/task';

interface KanbanCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: TaskStatus) => void;
}

export const KanbanCard = ({ task, onEdit, onDelete, onStatusChange }: KanbanCardProps) => {
  return (
    <div
      style={{
        backgroundColor: '#ffffff',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        padding: '12px',
        cursor: 'pointer',
        transition: 'all 0.2s',
        boxShadow: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.boxShadow = '0 4px 6px -1px rgb(0 0 0 / 0.1)';
        e.currentTarget.style.transform = 'translateY(-2px)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.boxShadow = '0 1px 2px 0 rgb(0 0 0 / 0.05)';
        e.currentTarget.style.transform = 'translateY(0)';
      }}
    >
      <div style={{ marginBottom: '8px' }}>
        <h4
          style={{
            margin: 0,
            fontSize: '15px',
            fontWeight: '600',
            color: '#1f2937',
            lineHeight: '1.4',
          }}
        >
          {task.title}
        </h4>
      </div>

      {task.description && (
        <p
          style={{
            margin: '0 0 12px 0',
            fontSize: '13px',
            color: '#6b7280',
            lineHeight: '1.5',
            display: '-webkit-box',
            WebkitLineClamp: 3,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
          }}
        >
          {task.description}
        </p>
      )}

      {task.due_at && (
        <div
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: '4px',
            padding: '4px 8px',
            backgroundColor: '#fef3c7',
            borderRadius: '4px',
            marginBottom: '12px',
          }}
        >
          <span style={{ fontSize: '12px', color: '#92400e' }}>
            📅 {format(new Date(task.due_at), 'MMM d, yyyy')}
          </span>
        </div>
      )}

      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginTop: '12px',
          paddingTop: '12px',
          borderTop: '1px solid #f3f4f6',
        }}
      >
        <select
          value={task.status}
          onChange={(e) => {
            e.stopPropagation();
            onStatusChange(task.id, e.target.value as TaskStatus);
          }}
          onClick={(e) => e.stopPropagation()}
          style={{
            padding: '4px 8px',
            fontSize: '12px',
            fontWeight: '500',
            border: '1px solid #d1d5db',
            borderRadius: '4px',
            backgroundColor: '#ffffff',
            color: '#374151',
            cursor: 'pointer',
          }}
        >
          <option value="todo">To Do</option>
          <option value="in_progress">In Progress</option>
          <option value="done">Done</option>
        </select>

        <div style={{ display: 'flex', gap: '6px' }}>
          <button
            onClick={(e) => {
              e.stopPropagation();
              onEdit(task);
            }}
            style={{
              padding: '4px 10px',
              fontSize: '12px',
              fontWeight: '500',
              border: '1px solid #d1d5db',
              borderRadius: '4px',
              backgroundColor: '#ffffff',
              color: '#374151',
              cursor: 'pointer',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#f9fafb';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#ffffff';
            }}
          >
            ✏️
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              onDelete(task.id);
            }}
            style={{
              padding: '4px 10px',
              fontSize: '12px',
              fontWeight: '600',
              border: '1px solid #dc2626',
              borderRadius: '4px',
              backgroundColor: '#dc2626',
              color: '#ffffff',
              cursor: 'pointer',
              transition: 'all 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#b91c1c';
              e.currentTarget.style.transform = 'scale(1.05)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#dc2626';
              e.currentTarget.style.transform = 'scale(1)';
            }}
            title="タスクを削除"
          >
            削除
          </button>
        </div>
      </div>
    </div>
  );
};
