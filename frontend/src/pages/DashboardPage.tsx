import { useAuth } from '../hooks/useAuth';
import { TaskList } from '../components/Tasks/TaskList';
import './DashboardPage.css';

export function DashboardPage() {
  const { user, logout } = useAuth();

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
        <TaskList />
      </div>
    </div>
  );
}
