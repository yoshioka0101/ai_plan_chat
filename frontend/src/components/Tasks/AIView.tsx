import { useState, useRef, useEffect } from 'react';
import { interpretationService } from '../../services/interpretationService';
import { interpretationItemService } from '../../services/interpretationItemService';
import type {
  AIInterpretation,
  InterpretationResponse,
  InterpretationItem,
  InterpretationItemData,
} from '../../types/interpretation';

interface Message {
  id: string;
  type: 'user' | 'ai';
  content: string;
  interpretation?: AIInterpretation;
  timestamp: Date;
}

interface AIViewProps {
  onNewTaskClick: () => void;
  onTaskCreated?: (taskId: string) => Promise<void> | void;
  onNotify?: (message: string) => void;
}

export function AIView({ onNewTaskClick, onTaskCreated, onNotify }: AIViewProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputText, setInputText] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeInterpretationId, setActiveInterpretationId] = useState<string | null>(null);
  const [items, setItems] = useState<InterpretationItem[]>([]);
  const [itemsLoading, setItemsLoading] = useState(false);
  const [itemEdits, setItemEdits] = useState<Record<string, InterpretationItemData>>({});
  const [approvingItemId, setApprovingItemId] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const loadItems = async (interpretationId: string) => {
    try {
      setItemsLoading(true);
      const data = await interpretationItemService.getItems(interpretationId);
      setItems(data);
      setItemEdits(
        data.reduce<Record<string, InterpretationItemData>>((acc, item) => {
          acc[item.id] = { ...item.data };
          return acc;
        }, {})
      );
    } catch (err) {
      setError('アイテムの取得に失敗しました。もう一度お試しください。');
      console.error('Error fetching interpretation items:', err);
    } finally {
      setItemsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputText.trim() || isLoading) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: 'user',
      content: inputText,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputText('');
    setIsLoading(true);
    setError(null);

    try {
      const response: InterpretationResponse =
        await interpretationService.createInterpretation({
          input_text: inputText,
        });

      const interpretationId = response.interpretation.id;
      setActiveInterpretationId(interpretationId);
      loadItems(interpretationId);

      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: 'ai',
        content: formatInterpretationResponse(response),
        interpretation: response.interpretation,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, aiMessage]);
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || err.message || 'Failed to get AI interpretation';
      const errorDetails = err.response?.data?.error || '';
      setError(`${errorMessage}${errorDetails ? ': ' + errorDetails : ''}. Please try again.`);
      console.error('Error creating interpretation:', err);
      console.error('Error response:', err.response?.data);
      console.error('Error status:', err.response?.status);
    } finally {
      setIsLoading(false);
    }
  };

  const formatInterpretationResponse = (response: InterpretationResponse): string => {
    const { interpretation } = response;
    const result = interpretation.structured_result;

    let message = `I've analyzed your input:\n\n`;

    if (result.title) {
      message += `**Title:** ${result.title}\n\n`;
    }

    if (result.description) {
      message += `**Description:** ${result.description}\n\n`;
    }

    if (result.type) {
      message += `**Type:** ${result.type}\n\n`;
    }

    if (result.metadata) {
      message += `**Details:**\n`;
      if (result.metadata.priority) {
        message += `- Priority: ${result.metadata.priority}\n`;
      }
      if (result.metadata.deadline) {
        message += `- Deadline: ${new Date(result.metadata.deadline).toLocaleDateString()}\n`;
      }
      if (result.metadata.tags && result.metadata.tags.length > 0) {
        message += `- Tags: ${result.metadata.tags.join(', ')}\n`;
      }
    }

    return message;
  };

  const handleFieldChange = (itemId: string, field: keyof InterpretationItemData, value: string | string[]) => {
    setItemEdits((prev) => {
      const next = { ...prev[itemId] };

      if ((typeof value === 'string' && value === '') || (Array.isArray(value) && value.length === 0)) {
        delete next[field];
      } else {
        next[field] = value;
      }

      return {
        ...prev,
        [itemId]: next,
      };
    });
  };

  // 編集内容は承認時に反映されるため、保存ボタンは撤去

  const handleApproveItem = async (itemId: string) => {
    try {
      setApprovingItemId(itemId);
      const response = await interpretationItemService.approveItem(itemId);
      setItems((prev) =>
        prev.map((item) =>
          item.id === itemId
            ? { ...item, status: 'created', resource_id: response.resource_id, reviewed_at: new Date().toISOString() }
            : item
        )
      );

      if (onTaskCreated) {
        await onTaskCreated(response.resource_id);
      }
      if (onNotify) {
        onNotify('タスクを作成しました');
      }
    } catch (err) {
      setError('承認に失敗しました。もう一度お試しください。');
      console.error('Error approving item:', err);
    } finally {
      setApprovingItemId(null);
    }
  };

  const formatMessageContent = (content: string) => {
    return content.split('\n').map((line, index) => {
      if (line.startsWith('**') && line.endsWith('**')) {
        const text = line.slice(2, -2);
        return (
          <strong key={index} style={{ display: 'block', marginTop: index > 0 ? '8px' : '0' }}>
            {text}
          </strong>
        );
      }
      return (
        <span key={index} style={{ display: 'block' }}>
          {line || '\u00A0'}
        </span>
      );
    });
  };

  return (
    <div style={{ width: '100%', height: 'calc(100vh - 200px)', display: 'flex', flexDirection: 'column' }}>
      {/* Header with New Task Button */}
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '24px',
        padding: '16px 24px',
        backgroundColor: '#ffffff',
        borderRadius: '12px',
        boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
      }}>
        <div>
          <h2 style={{ margin: '0 0 4px 0', fontSize: '20px', fontWeight: '700', color: '#1f2937' }}>
            🤖 AI Assistant
          </h2>
          <p style={{ margin: 0, fontSize: '14px', color: '#6b7280' }}>
            Tell me what you need to do, and I'll help you create tasks
          </p>
        </div>
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

      {/* Chat Container */}
      <div style={{
        flex: 1,
        backgroundColor: '#ffffff',
        borderRadius: '12px',
        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
      }}>
        {/* Messages */}
        <div style={{
          flex: 1,
          overflowY: 'auto',
          padding: '24px',
          display: 'flex',
          flexDirection: 'column',
          gap: '16px',
        }}>
          {messages.length === 0 && (
            <div style={{ textAlign: 'center', padding: '40px 20px', color: '#6b7280' }}>
              <h3 style={{ margin: '0 0 16px 0', color: '#1f2937' }}>Welcome!</h3>
              <p style={{ margin: '0 0 16px 0' }}>Try asking me things like:</p>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', maxWidth: '600px', margin: '0 auto' }}>
                <div style={{
                  padding: '12px 16px',
                  backgroundColor: '#f9fafb',
                  borderRadius: '8px',
                  fontStyle: 'italic',
                  textAlign: 'left',
                }}>
                  "明日、美容院に行く予定を立ててください"
                </div>
                <div style={{
                  padding: '12px 16px',
                  backgroundColor: '#f9fafb',
                  borderRadius: '8px',
                  fontStyle: 'italic',
                  textAlign: 'left',
                }}>
                  "Remind me to call John tomorrow at 3pm"
                </div>
                <div style={{
                  padding: '12px 16px',
                  backgroundColor: '#f9fafb',
                  borderRadius: '8px',
                  fontStyle: 'italic',
                  textAlign: 'left',
                }}>
                  "Create a high priority task to review the code"
                </div>
              </div>
            </div>
          )}

          {messages.map((message) => (
            <div
              key={message.id}
              style={{
                display: 'flex',
                flexDirection: 'column',
                alignSelf: message.type === 'user' ? 'flex-end' : 'flex-start',
                maxWidth: '70%',
              }}
            >
              <div
                style={{
                  padding: '12px 16px',
                  borderRadius: '12px',
                  backgroundColor: message.type === 'user' ? '#3b82f6' : '#f3f4f6',
                  color: message.type === 'user' ? '#ffffff' : '#1f2937',
                }}
              >
                <div style={{ fontSize: '12px', fontWeight: '600', marginBottom: '4px', opacity: 0.8 }}>
                  {message.type === 'user' ? 'You' : 'AI Assistant'}
                </div>
                <div style={{ lineHeight: '1.5', whiteSpace: 'pre-wrap', wordWrap: 'break-word' }}>
                  {formatMessageContent(message.content)}
                </div>
              </div>
              <div style={{
                fontSize: '11px',
                color: '#9ca3af',
                marginTop: '4px',
                paddingLeft: message.type === 'user' ? '0' : '8px',
                paddingRight: message.type === 'user' ? '8px' : '0',
                textAlign: message.type === 'user' ? 'right' : 'left',
              }}>
                {message.timestamp.toLocaleTimeString()}
              </div>
            </div>
          ))}

          {isLoading && (
            <div style={{
              alignSelf: 'flex-start',
              maxWidth: '70%',
              padding: '12px 16px',
              borderRadius: '12px',
              backgroundColor: '#f3f4f6',
              color: '#6b7280',
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <span>Thinking</span>
                <span className="loading-dots">...</span>
              </div>
            </div>
          )}

          {error && (
            <div style={{
              padding: '12px 16px',
              backgroundColor: '#fee2e2',
              color: '#991b1b',
              borderRadius: '8px',
              textAlign: 'center',
            }}>
              {error}
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        {/* Input Form */}
        <form
          onSubmit={handleSubmit}
          style={{
            padding: '16px 24px',
            borderTop: '1px solid #e5e7eb',
            backgroundColor: '#ffffff',
            display: 'flex',
            gap: '12px',
          }}
        >
          <input
            type="text"
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            placeholder="Type your message here..."
            disabled={isLoading}
            style={{
              flex: 1,
              padding: '12px 16px',
              border: '1px solid #d1d5db',
              borderRadius: '8px',
              fontSize: '14px',
              outline: 'none',
              transition: 'border-color 0.2s',
            }}
            onFocus={(e) => {
              e.currentTarget.style.borderColor = '#3b82f6';
              e.currentTarget.style.boxShadow = '0 0 0 3px rgba(59, 130, 246, 0.1)';
            }}
            onBlur={(e) => {
              e.currentTarget.style.borderColor = '#d1d5db';
              e.currentTarget.style.boxShadow = 'none';
            }}
          />
          <button
            type="submit"
            disabled={isLoading || !inputText.trim()}
            style={{
              padding: '12px 24px',
              backgroundColor: isLoading || !inputText.trim() ? '#9ca3af' : '#3b82f6',
              color: '#ffffff',
              border: 'none',
              borderRadius: '8px',
              fontWeight: '600',
              cursor: isLoading || !inputText.trim() ? 'not-allowed' : 'pointer',
              transition: 'background-color 0.2s',
            }}
            onMouseEnter={(e) => {
              if (!isLoading && inputText.trim()) {
                e.currentTarget.style.backgroundColor = '#2563eb';
              }
            }}
            onMouseLeave={(e) => {
              if (!isLoading && inputText.trim()) {
                e.currentTarget.style.backgroundColor = '#3b82f6';
              }
            }}
          >
            Send
          </button>
        </form>
      </div>

      {/* Interpretation items */}
      {activeInterpretationId && (
        <div
          style={{
            marginTop: '24px',
            backgroundColor: '#ffffff',
            borderRadius: '12px',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
            padding: '20px',
          }}
        >
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', marginBottom: '12px' }}>
            <h3 style={{ margin: 0, fontSize: '18px', fontWeight: 700, color: '#1f2937' }}>Review AI suggestions</h3>
            <p style={{ margin: 0, color: '#6b7280', fontSize: '13px' }}>編集して承認するとタスクが作成されます</p>
          </div>

          {itemsLoading ? (
            <div style={{ color: '#6b7280', fontSize: '14px' }}>読み込み中...</div>
          ) : items.length === 0 ? (
            <div style={{ color: '#6b7280', fontSize: '14px' }}>提案されたアイテムがありません。</div>
          ) : (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '16px' }}>
              {items.map((item) => {
                const draft = itemEdits[item.id] || {};
                const tagsText = draft.tags?.join(', ') || '';
                const dueAt = draft.due_at ? new Date(draft.due_at).toISOString().slice(0, 16) : '';

                const isPending = item.status === 'pending';

                return (
                  <div
                    key={item.id}
                    style={{
                      border: '1px solid #e5e7eb',
                      borderRadius: '12px',
                      padding: '16px',
                      display: 'flex',
                      flexDirection: 'column',
                      gap: '10px',
                    }}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <span style={{ fontWeight: 700, color: '#111827' }}>Task #{item.item_index + 1}</span>
                      <span
                        style={{
                          padding: '4px 10px',
                          borderRadius: '999px',
                          fontSize: '12px',
                          backgroundColor: isPending ? '#fff7ed' : '#ecfdf3',
                          color: isPending ? '#c2410c' : '#166534',
                          border: `1px solid ${isPending ? '#fed7aa' : '#bbf7d0'}`,
                        }}
                      >
                        {isPending ? 'pending' : 'created'}
                      </span>
                    </div>

                    <label style={{ display: 'flex', flexDirection: 'column', gap: '6px', fontSize: '13px', color: '#374151' }}>
                      タイトル
                      <input
                        type="text"
                        value={draft.title || ''}
                        onChange={(e) => handleFieldChange(item.id, 'title', e.target.value)}
                        disabled={!isPending}
                        style={{
                          padding: '10px',
                          borderRadius: '8px',
                          border: '1px solid #d1d5db',
                        }}
                      />
                    </label>

                    <label style={{ display: 'flex', flexDirection: 'column', gap: '6px', fontSize: '13px', color: '#374151' }}>
                      詳細
                      <textarea
                        value={draft.description || ''}
                        onChange={(e) => handleFieldChange(item.id, 'description', e.target.value)}
                        disabled={!isPending}
                        rows={3}
                        style={{
                          padding: '10px',
                          borderRadius: '8px',
                          border: '1px solid #d1d5db',
                          resize: 'vertical',
                        }}
                      />
                    </label>

                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
                      <label style={{ display: 'flex', flexDirection: 'column', gap: '6px', fontSize: '13px', color: '#374151' }}>
                        期限
                        <input
                          type="datetime-local"
                          value={dueAt}
                          onChange={(e) => handleFieldChange(item.id, 'due_at', e.target.value ? new Date(e.target.value).toISOString() : '')}
                          disabled={!isPending}
                          style={{
                            padding: '10px',
                            borderRadius: '8px',
                            border: '1px solid #d1d5db',
                          }}
                        />
                      </label>
                      <label style={{ display: 'flex', flexDirection: 'column', gap: '6px', fontSize: '13px', color: '#374151' }}>
                        ステータス
                        <select
                          value={draft.status || ''}
                          onChange={(e) => handleFieldChange(item.id, 'status', e.target.value)}
                          disabled={!isPending}
                          style={{
                            padding: '10px',
                            borderRadius: '8px',
                            border: '1px solid #d1d5db',
                          }}
                        >
                          <option value="">未指定</option>
                          <option value="todo">To Do</option>
                          <option value="in_progress">In Progress</option>
                          <option value="done">Done</option>
                        </select>
                      </label>
                    </div>

                    <label style={{ display: 'flex', flexDirection: 'column', gap: '6px', fontSize: '13px', color: '#374151' }}>
                      タグ (カンマ区切り)
                      <input
                        type="text"
                        value={tagsText}
                        onChange={(e) =>
                          handleFieldChange(
                            item.id,
                            'tags',
                            e.target.value
                              .split(',')
                              .map((t) => t.trim())
                              .filter(Boolean)
                          )
                        }
                        disabled={!isPending}
                        style={{
                          padding: '10px',
                          borderRadius: '8px',
                          border: '1px solid #d1d5db',
                        }}
                      />
                    </label>

                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px', marginTop: '4px' }}>
                      <button
                        onClick={() => handleApproveItem(item.id)}
                        disabled={!isPending || approvingItemId === item.id}
                        style={{
                          minWidth: '160px',
                          padding: '10px 14px',
                          borderRadius: '8px',
                          border: 'none',
                          backgroundColor: isPending ? '#10b981' : '#9ca3af',
                          color: '#fff',
                          cursor: isPending && approvingItemId !== item.id ? 'pointer' : 'not-allowed',
                        }}
                      >
                        {approvingItemId === item.id ? '承認中...' : '承認して作成'}
                      </button>
                    </div>

                    {item.resource_id && (
                      <div style={{ fontSize: '12px', color: '#16a34a' }}>
                        作成済みリソースID: {item.resource_id}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
