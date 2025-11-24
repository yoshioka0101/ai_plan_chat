import { useState, useRef, useEffect } from 'react';
import { interpretationService } from '../../services/interpretationService';
import type {
  AIInterpretation,
  InterpretationResponse,
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
}

export function AIView({ onNewTaskClick }: AIViewProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputText, setInputText] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

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
                  "I need to finish the project report by Friday"
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
    </div>
  );
}
