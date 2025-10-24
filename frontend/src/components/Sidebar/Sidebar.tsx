type ViewMode = 'kanban' | 'list' | 'calendar';

interface SidebarProps {
  onNewTaskClick: () => void;
  isOpen: boolean;
  onToggle: (isOpen: boolean) => void;
  viewMode: ViewMode;
  onViewChange: (mode: ViewMode) => void;
}

export const Sidebar = ({ onNewTaskClick, isOpen, onToggle, viewMode, onViewChange }: SidebarProps) => {

  return (
    <>
      {/* ハンバーガーボタン */}
      <button
        onClick={() => onToggle(!isOpen)}
        style={{
          position: 'fixed',
          top: '20px',
          left: isOpen ? '320px' : '20px',
          zIndex: 1000,
          padding: '12px',
          backgroundColor: '#ffffff',
          border: '2px solid #e5e7eb',
          borderRadius: '8px',
          cursor: 'pointer',
          boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
          transition: 'left 0.3s ease',
          display: 'flex',
          flexDirection: 'column',
          gap: '4px',
          width: '40px',
          height: '40px',
          alignItems: 'center',
          justifyContent: 'center',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.backgroundColor = '#f9fafb';
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.backgroundColor = '#ffffff';
        }}
      >
        <div
          style={{
            width: '20px',
            height: '2px',
            backgroundColor: '#374151',
            transition: 'all 0.3s',
            transform: isOpen ? 'rotate(45deg) translateY(6px)' : 'none',
          }}
        />
        <div
          style={{
            width: '20px',
            height: '2px',
            backgroundColor: '#374151',
            transition: 'all 0.3s',
            opacity: isOpen ? 0 : 1,
          }}
        />
        <div
          style={{
            width: '20px',
            height: '2px',
            backgroundColor: '#374151',
            transition: 'all 0.3s',
            transform: isOpen ? 'rotate(-45deg) translateY(-6px)' : 'none',
          }}
        />
      </button>

      {/* サイドバー */}
      <div
        style={{
          position: 'fixed',
          top: 0,
          left: 0,
          bottom: 0,
          width: '320px',
          backgroundColor: '#ffffff',
          borderRight: '1px solid #e5e7eb',
          boxShadow: '2px 0 8px rgb(0 0 0 / 0.1)',
          transform: isOpen ? 'translateX(0)' : 'translateX(-100%)',
          transition: 'transform 0.3s ease',
          zIndex: 999,
          overflow: 'auto',
        }}
      >
        <div style={{ padding: '80px 20px 20px 20px' }}>
          {/* サイドバーヘッダー */}
          <div style={{ marginBottom: '24px' }}>
            <h2
              style={{
                margin: '0 0 8px 0',
                fontSize: '24px',
                fontWeight: '700',
                color: '#1f2937',
              }}
            >
              📋 Task Manager
            </h2>
            <p style={{ margin: 0, fontSize: '14px', color: '#6b7280' }}>
              Organize your tasks efficiently
            </p>
          </div>

          {/* タスク作成ボタン */}
          <button
            onClick={onNewTaskClick}
            style={{
              width: '100%',
              padding: '14px 20px',
              borderRadius: '8px',
              border: 'none',
              backgroundColor: '#3b82f6',
              color: '#ffffff',
              cursor: 'pointer',
              fontSize: '15px',
              fontWeight: '600',
              boxShadow: '0 2px 4px rgb(0 0 0 / 0.1)',
              marginBottom: '24px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '8px',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#2563eb';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#3b82f6';
            }}
          >
            <span style={{ fontSize: '18px' }}>+</span>
            新規タスク作成
          </button>

          {/* メニュー項目 */}
          <div style={{ marginTop: '32px' }}>
            <div
              onClick={() => onViewChange('kanban')}
              style={{
                padding: '12px 16px',
                marginBottom: '8px',
                borderRadius: '8px',
                cursor: 'pointer',
                transition: 'background-color 0.2s',
                backgroundColor: viewMode === 'kanban' ? '#f3f4f6' : 'transparent',
                fontWeight: viewMode === 'kanban' ? '600' : '400',
                color: viewMode === 'kanban' ? '#1f2937' : '#6b7280',
              }}
              onMouseEnter={(e) => {
                if (viewMode !== 'kanban') {
                  e.currentTarget.style.backgroundColor = '#f9fafb';
                }
              }}
              onMouseLeave={(e) => {
                if (viewMode !== 'kanban') {
                  e.currentTarget.style.backgroundColor = 'transparent';
                }
              }}
            >
              📊 Kanban Board
            </div>
            <div
              onClick={() => onViewChange('list')}
              style={{
                padding: '12px 16px',
                marginBottom: '8px',
                borderRadius: '8px',
                cursor: 'pointer',
                transition: 'background-color 0.2s',
                backgroundColor: viewMode === 'list' ? '#f3f4f6' : 'transparent',
                fontWeight: viewMode === 'list' ? '600' : '400',
                color: viewMode === 'list' ? '#1f2937' : '#6b7280',
              }}
              onMouseEnter={(e) => {
                if (viewMode !== 'list') {
                  e.currentTarget.style.backgroundColor = '#f9fafb';
                }
              }}
              onMouseLeave={(e) => {
                if (viewMode !== 'list') {
                  e.currentTarget.style.backgroundColor = 'transparent';
                }
              }}
            >
              📝 List View
            </div>
            <div
              onClick={() => onViewChange('calendar')}
              style={{
                padding: '12px 16px',
                marginBottom: '8px',
                borderRadius: '8px',
                cursor: 'pointer',
                transition: 'background-color 0.2s',
                backgroundColor: viewMode === 'calendar' ? '#f3f4f6' : 'transparent',
                fontWeight: viewMode === 'calendar' ? '600' : '400',
                color: viewMode === 'calendar' ? '#1f2937' : '#6b7280',
              }}
              onMouseEnter={(e) => {
                if (viewMode !== 'calendar') {
                  e.currentTarget.style.backgroundColor = '#f9fafb';
                }
              }}
              onMouseLeave={(e) => {
                if (viewMode !== 'calendar') {
                  e.currentTarget.style.backgroundColor = 'transparent';
                }
              }}
            >
              📅 Calendar
            </div>
          </div>
        </div>
      </div>

      {/* オーバーレイ（モバイル用） */}
      {isOpen && (
        <div
          onClick={() => onToggle(false)}
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.3)',
            zIndex: 998,
            display: window.innerWidth < 768 ? 'block' : 'none',
          }}
        />
      )}
    </>
  );
};
