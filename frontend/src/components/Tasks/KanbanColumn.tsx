import type { Task, TaskStatus } from '../../types/task';
import { KanbanCard } from './KanbanCard';

interface KanbanColumnProps {
  title: string;
  status: TaskStatus;
  tasks: Task[];
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: TaskStatus) => void;
}

const getColumnColor = (status: TaskStatus) => {
  switch (status) {
    case 'todo':
      return '#fbbf24'; // amber
    case 'in_progress':
      return '#60a5fa'; // blue
    case 'done':
      return '#4ade80'; // green
    default:
      return '#94a3b8'; // gray
  }
};

export const KanbanColumn = ({
  title,
  status,
  tasks,
  onEdit,
  onDelete,
  onStatusChange
}: KanbanColumnProps) => {
  const columnColor = getColumnColor(status);

  return (
    <div
      style={{
        backgroundColor: '#ffffff',
        borderRadius: '12px',
        padding: '20px',
        minHeight: '600px',
        display: 'flex',
        flexDirection: 'column',
        width: '100%',
        border: '1px solid #e5e7eb',
        boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.1)',
      }}
    >
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '12px',
          marginBottom: '20px',
          paddingBottom: '16px',
          borderBottom: `3px solid ${columnColor}`,
        }}
      >
        <h3
          style={{
            margin: 0,
            fontSize: '18px',
            fontWeight: '700',
            color: '#1f2937',
            flex: 1,
          }}
        >
          {title}
        </h3>
        <span
          style={{
            backgroundColor: columnColor + '20',
            color: columnColor,
            padding: '6px 14px',
            borderRadius: '16px',
            fontSize: '14px',
            fontWeight: '700',
            minWidth: '32px',
            textAlign: 'center',
          }}
        >
          {tasks.length}
        </span>
      </div>

      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: '12px',
          flex: 1,
        }}
      >
        {tasks.length === 0 ? (
          <div
            style={{
              textAlign: 'center',
              padding: '40px 20px',
              color: '#9ca3af',
              fontSize: '14px',
            }}
          >
            No tasks
          </div>
        ) : (
          tasks.map((task) => (
            <KanbanCard
              key={task.id}
              task={task}
              onEdit={onEdit}
              onDelete={onDelete}
              onStatusChange={onStatusChange}
            />
          ))
        )}
      </div>
    </div>
  );
};
