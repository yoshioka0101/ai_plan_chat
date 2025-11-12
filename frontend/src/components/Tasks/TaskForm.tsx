import { useState, useEffect } from 'react';
import type { Task, CreateTaskRequest, TaskStatus } from '../../types/task';
import { interpretationService } from '../../services/interpretationService';

interface TaskFormProps {
  task?: Task;
  onSubmit: (task: CreateTaskRequest) => void;
  onCancel: () => void;
}

export const TaskForm = ({ task, onSubmit, onCancel }: TaskFormProps) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [dueDate, setDueDate] = useState('');
  const [status, setStatus] = useState<TaskStatus>('todo');
  const [aiInput, setAiInput] = useState('');
  const [isAILoading, setIsAILoading] = useState(false);
  const [showAIInput, setShowAIInput] = useState(false);

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description || '');
      setDueDate(task.due_at ? task.due_at.split('T')[0] : '');
      setStatus(task.status);
    } else {
      // 新規作成時はフォームをリセット
      setTitle('');
      setDescription('');
      setDueDate('');
      setStatus('todo');
    }
  }, [task]);

  const handleAIGenerate = async () => {
    if (!aiInput.trim()) return;

    setIsAILoading(true);
    try {
      const response = await interpretationService.createInterpretation({
        input_text: aiInput,
      });

      const result = response.interpretation.structured_result;

      // AI結果をフォームに反映
      if (result.title) {
        setTitle(result.title);
      }
      if (result.description) {
        setDescription(result.description);
      }
      if (result.metadata?.deadline) {
        const deadlineDate = new Date(result.metadata.deadline);
        setDueDate(deadlineDate.toISOString().split('T')[0]);
      }
      // priority から status への変換（オプション）
      if (result.metadata?.priority === 'high') {
        setStatus('in_progress');
      }

      setShowAIInput(false);
      setAiInput('');
    } catch (error) {
      console.error('AI generation failed:', error);
      alert('AI generation failed. Please try again.');
    } finally {
      setIsAILoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const taskData: CreateTaskRequest = {
      title,
      description: description || undefined,
      due_at: dueDate || undefined,
      status: status,
    };

    onSubmit(taskData);
  };

  return (
    <div
      style={{
        backgroundColor: '#f9fafb',
        padding: '20px',
        borderRadius: '12px',
        border: '1px solid #e5e7eb',
      }}
    >
      <h3 style={{ margin: '0 0 16px 0', fontSize: '16px', fontWeight: '700', color: '#1f2937' }}>
        {task ? '✏️ タスクを編集' : '✨ タスクを作成'}
      </h3>

      {!task && !showAIInput && (
        <div style={{ marginBottom: '16px' }}>
          <button
            type="button"
            onClick={() => setShowAIInput(true)}
            style={{
              width: '100%',
              padding: '12px',
              borderRadius: '8px',
              border: '2px dashed #3b82f6',
              backgroundColor: '#eff6ff',
              color: '#3b82f6',
              cursor: 'pointer',
              fontSize: '14px',
              fontWeight: '600',
              transition: 'all 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#dbeafe';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#eff6ff';
            }}
          >
            🤖 AIでタスクを生成
          </button>
        </div>
      )}

      {showAIInput && (
        <div style={{ marginBottom: '16px', padding: '16px', backgroundColor: '#eff6ff', borderRadius: '8px' }}>
          <label style={{ display: 'block', marginBottom: '8px', fontWeight: '600', fontSize: '14px', color: '#1e40af' }}>
            AIに何をすべきか教えてください
          </label>
          <textarea
            value={aiInput}
            onChange={(e) => setAiInput(e.target.value)}
            placeholder="例: 明日までにプロジェクトレポートを完成させる"
            rows={3}
            disabled={isAILoading}
            style={{
              width: '100%',
              padding: '10px',
              border: '1px solid #3b82f6',
              borderRadius: '6px',
              fontSize: '14px',
              fontFamily: 'inherit',
              boxSizing: 'border-box',
              resize: 'vertical',
              marginBottom: '8px',
            }}
          />
          <div style={{ display: 'flex', gap: '8px' }}>
            <button
              type="button"
              onClick={handleAIGenerate}
              disabled={isAILoading || !aiInput.trim()}
              style={{
                flex: 1,
                padding: '10px',
                borderRadius: '6px',
                border: 'none',
                backgroundColor: isAILoading || !aiInput.trim() ? '#9ca3af' : '#3b82f6',
                color: '#ffffff',
                cursor: isAILoading || !aiInput.trim() ? 'not-allowed' : 'pointer',
                fontSize: '14px',
                fontWeight: '600',
              }}
            >
              {isAILoading ? '生成中...' : '✨ 生成'}
            </button>
            <button
              type="button"
              onClick={() => {
                setShowAIInput(false);
                setAiInput('');
              }}
              disabled={isAILoading}
              style={{
                padding: '10px 16px',
                borderRadius: '6px',
                border: '1px solid #d1d5db',
                backgroundColor: '#ffffff',
                color: '#374151',
                cursor: isAILoading ? 'not-allowed' : 'pointer',
                fontSize: '14px',
                fontWeight: '600',
              }}
            >
              キャンセル
            </button>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '16px' }}>
          <label style={{ display: 'block', marginBottom: '6px', fontWeight: '500', fontSize: '14px' }}>
            Title *
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            placeholder="Enter task title"
            style={{
              width: '100%',
              padding: '10px',
              border: '1px solid #d1d5db',
              borderRadius: '6px',
              fontSize: '14px',
              boxSizing: 'border-box',
            }}
          />
        </div>

        <div style={{ marginBottom: '16px' }}>
          <label style={{ display: 'block', marginBottom: '6px', fontWeight: '500', fontSize: '14px' }}>
            Description
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Enter task description (optional)"
            rows={3}
            style={{
              width: '100%',
              padding: '10px',
              border: '1px solid #d1d5db',
              borderRadius: '6px',
              fontSize: '14px',
              fontFamily: 'inherit',
              boxSizing: 'border-box',
              resize: 'vertical',
            }}
          />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginBottom: '16px' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '6px', fontWeight: '500', fontSize: '14px' }}>
              Due Date
            </label>
            <input
              type="date"
              value={dueDate}
              onChange={(e) => setDueDate(e.target.value)}
              style={{
                width: '100%',
                padding: '10px',
                border: '1px solid #d1d5db',
                borderRadius: '6px',
                fontSize: '14px',
                boxSizing: 'border-box',
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '6px', fontWeight: '500', fontSize: '14px' }}>
              Status *
            </label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value as TaskStatus)}
              required
              style={{
                width: '100%',
                padding: '10px',
                border: '1px solid #d1d5db',
                borderRadius: '6px',
                fontSize: '14px',
                boxSizing: 'border-box',
              }}
            >
              <option value="todo">To Do</option>
              <option value="in_progress">In Progress</option>
              <option value="done">Done</option>
            </select>
          </div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', marginTop: '16px' }}>
          <button
            type="submit"
            style={{
              width: '100%',
              padding: '12px',
              borderRadius: '8px',
              border: 'none',
              backgroundColor: '#3b82f6',
              color: '#ffffff',
              cursor: 'pointer',
              fontSize: '14px',
              fontWeight: '600',
              boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.1)',
              transition: 'all 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#2563eb';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#3b82f6';
            }}
          >
            {task ? 'タスクを更新' : 'タスクを作成'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            style={{
              width: '100%',
              padding: '12px',
              borderRadius: '8px',
              border: '1px solid #d1d5db',
              backgroundColor: '#ffffff',
              color: '#374151',
              cursor: 'pointer',
              fontSize: '14px',
              fontWeight: '600',
              transition: 'all 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#f9fafb';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#ffffff';
            }}
          >
            キャンセル
          </button>
        </div>
      </form>
    </div>
  );
};
