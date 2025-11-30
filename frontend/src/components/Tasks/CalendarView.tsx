import { useState } from 'react';
import { format, startOfMonth, endOfMonth, eachDayOfInterval, isSameMonth, isSameDay, addMonths, subMonths, startOfWeek, endOfWeek } from 'date-fns';
import { ja } from 'date-fns/locale';
import type { Task, TaskStatus } from '../../types/task';

interface CalendarViewProps {
  tasks: Task[];
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: TaskStatus) => void;
  onNewTaskClick: (date?: Date) => void;
}

export const CalendarView = ({ tasks, onEdit, onNewTaskClick }: CalendarViewProps) => {
  const [currentDate, setCurrentDate] = useState(new Date());

  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(currentDate);
  const calendarStart = startOfWeek(monthStart, { weekStartsOn: 0 }); // Sunday
  const calendarEnd = endOfWeek(monthEnd, { weekStartsOn: 0 });

  const calendarDays = eachDayOfInterval({ start: calendarStart, end: calendarEnd });

  const getTasksForDate = (date: Date) => {
    return tasks.filter(task => {
      if (!task.due_at) return false;
      return isSameDay(new Date(task.due_at), date);
    });
  };

  const getStatusColor = (status: TaskStatus) => {
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

  const handlePrevMonth = () => {
    setCurrentDate(subMonths(currentDate, 1));
  };

  const handleNextMonth = () => {
    setCurrentDate(addMonths(currentDate, 1));
  };

  const handleToday = () => {
    setCurrentDate(new Date());
  };

  const weekDays = ['日', '月', '火', '水', '木', '金', '土'];

  return (
    <div style={{ width: '100%' }}>
      {/* ヘッダー */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h2 style={{ margin: 0, fontSize: '24px', fontWeight: '700', color: '#1f2937' }}>
          {format(currentDate, 'yyyy年 M月', { locale: ja })}
        </h2>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={handlePrevMonth}
            style={{
              padding: '8px 16px',
              fontSize: '14px',
              fontWeight: '600',
              border: '1px solid #d1d5db',
              borderRadius: '6px',
              backgroundColor: '#ffffff',
              color: '#374151',
              cursor: 'pointer',
            }}
          >
            ← 前月
          </button>
          <button
            onClick={handleToday}
            style={{
              padding: '8px 16px',
              fontSize: '14px',
              fontWeight: '600',
              border: '1px solid #d1d5db',
              borderRadius: '6px',
              backgroundColor: '#ffffff',
              color: '#374151',
              cursor: 'pointer',
            }}
          >
            今日
          </button>
          <button
            onClick={handleNextMonth}
            style={{
              padding: '8px 16px',
              fontSize: '14px',
              fontWeight: '600',
              border: '1px solid #d1d5db',
              borderRadius: '6px',
              backgroundColor: '#ffffff',
              color: '#374151',
              cursor: 'pointer',
            }}
          >
            次月 →
          </button>
        </div>
      </div>

      {/* カレンダーグリッド */}
      <div
        style={{
          backgroundColor: '#ffffff',
          borderRadius: '12px',
          border: '1px solid #e5e7eb',
          overflow: 'hidden',
        }}
      >
        {/* 曜日ヘッダー */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(7, 1fr)',
            backgroundColor: '#f9fafb',
            borderBottom: '1px solid #e5e7eb',
          }}
        >
          {weekDays.map((day, index) => (
            <div
              key={day}
              style={{
                padding: '12px',
                textAlign: 'center',
                fontSize: '14px',
                fontWeight: '600',
                color: index === 0 ? '#dc2626' : index === 6 ? '#2563eb' : '#374151',
              }}
            >
              {day}
            </div>
          ))}
        </div>

        {/* 日付グリッド */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(7, 1fr)',
            gridAutoRows: 'minmax(120px, auto)',
          }}
        >
          {calendarDays.map((day, index) => {
            const dayTasks = getTasksForDate(day);
            const isCurrentMonth = isSameMonth(day, currentDate);
            const isToday = isSameDay(day, new Date());

            return (
              <div
                key={day.toISOString()}
                style={{
                  padding: '8px',
                  borderRight: (index + 1) % 7 !== 0 ? '1px solid #e5e7eb' : 'none',
                  borderBottom: index < calendarDays.length - 7 ? '1px solid #e5e7eb' : 'none',
                  backgroundColor: isToday ? '#fef3c7' : isCurrentMonth ? '#ffffff' : '#f9fafb',
                  cursor: 'pointer',
                  transition: 'background-color 0.2s',
                }}
                onMouseEnter={(e) => {
                  if (!isToday) {
                    e.currentTarget.style.backgroundColor = isCurrentMonth ? '#f9fafb' : '#f3f4f6';
                  }
                }}
                onMouseLeave={(e) => {
                  if (!isToday) {
                    e.currentTarget.style.backgroundColor = isCurrentMonth ? '#ffffff' : '#f9fafb';
                  }
                }}
                onClick={() => onNewTaskClick(day)}
              >
                {/* 日付 */}
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: '4px',
                  }}
                >
                  <span
                    style={{
                      fontSize: '14px',
                      fontWeight: isToday ? '700' : '600',
                      color: isCurrentMonth ? (isToday ? '#92400e' : '#1f2937') : '#9ca3af',
                    }}
                  >
                    {format(day, 'd')}
                  </span>
                  {dayTasks.length > 0 && (
                    <span
                      style={{
                        fontSize: '11px',
                        fontWeight: '600',
                        color: '#6b7280',
                        backgroundColor: '#e5e7eb',
                        padding: '2px 6px',
                        borderRadius: '10px',
                      }}
                    >
                      {dayTasks.length}
                    </span>
                  )}
                </div>

                {/* タスクリスト */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                  {dayTasks.slice(0, 3).map((task) => (
                    <div
                      key={task.id}
                      onClick={(e) => {
                        e.stopPropagation();
                        onEdit(task);
                      }}
                      style={{
                        padding: '4px 6px',
                        borderRadius: '4px',
                        fontSize: '12px',
                        fontWeight: '500',
                        backgroundColor: getStatusColor(task.status) + '20',
                        color: '#1f2937',
                        borderLeft: `3px solid ${getStatusColor(task.status)}`,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        cursor: 'pointer',
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.backgroundColor = getStatusColor(task.status) + '40';
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.backgroundColor = getStatusColor(task.status) + '20';
                      }}
                      title={task.title}
                    >
                      {task.title}
                    </div>
                  ))}
                  {dayTasks.length > 3 && (
                    <div
                      style={{
                        fontSize: '11px',
                        color: '#6b7280',
                        textAlign: 'center',
                        padding: '2px',
                      }}
                    >
                      +{dayTasks.length - 3} more
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
