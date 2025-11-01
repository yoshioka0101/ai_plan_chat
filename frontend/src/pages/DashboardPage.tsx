import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { TaskList } from '../components/Tasks/TaskList';
import './DashboardPage.css';

export function DashboardPage() {
  const { user, logout } = useAuth();
  const [activeTab, setActiveTab] = useState<'tasks' | 'ai-interpretation'>('tasks');

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div className="header-content">
          <h1 className="dashboard-title">AI Plan Chat</h1>
          <div className="user-menu">
            <span className="user-name">{user?.nickname || user?.email || 'User'}</span>
            <button onClick={logout} className="logout-button">
              Logout
            </button>
          </div>
        </div>
      </header>

      <div className="dashboard-content">
        <nav className="dashboard-nav">
          <button
            className={`nav-button ${activeTab === 'tasks' ? 'active' : ''}`}
            onClick={() => setActiveTab('tasks')}
          >
            Tasks
          </button>
          <button
            className={`nav-button ${activeTab === 'ai-interpretation' ? 'active' : ''}`}
            onClick={() => setActiveTab('ai-interpretation')}
          >
            AI Interpretation
          </button>
        </nav>

        <main className="dashboard-main">
          {activeTab === 'tasks' && (
            <div className="tab-content">
              <TaskList />
            </div>
          )}
          {activeTab === 'ai-interpretation' && (
            <div className="tab-content">
              <div className="coming-soon">
                <h2>AI Interpretation Feature</h2>
                <p>This feature is coming soon!</p>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  );
}
