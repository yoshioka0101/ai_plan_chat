import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export function AuthCallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { setAuth } = useAuth();
  const [error, setError] = useState<string>('');

  useEffect(() => {
    // Get token and user from URL params (sent by backend after OAuth)
    const token = searchParams.get('token');
    const userParam = searchParams.get('user');

    console.log('AuthCallback - token:', token ? 'exists' : 'missing');
    console.log('AuthCallback - userParam:', userParam ? 'exists' : 'missing');
    console.log('AuthCallback - full URL:', window.location.href);

    if (token && userParam) {
      try {
        const decodedUser = decodeURIComponent(userParam);
        console.log('AuthCallback - decoded user:', decodedUser);
        const user = JSON.parse(decodedUser);
        console.log('AuthCallback - parsed user:', user);
        setAuth(user, token);
        navigate('/dashboard', { replace: true });
      } catch (err) {
        console.error('Failed to parse user data:', err);
        setError('Authentication failed. Invalid user data.');
      }
    } else {
      console.error('AuthCallback - Missing params:', { token: !!token, userParam: !!userParam });
      setError('Authentication failed. Missing token or user data.');
    }
  }, [searchParams, setAuth, navigate]);

  if (error) {
    return (
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        backgroundColor: '#f3f4f6',
        padding: '2rem'
      }}>
        <div style={{
          backgroundColor: 'white',
          borderRadius: '8px',
          padding: '2rem',
          maxWidth: '400px',
          textAlign: 'center',
          boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.1)'
        }}>
          <h2 style={{ color: '#dc2626', marginBottom: '1rem' }}>Authentication Error</h2>
          <p style={{ color: '#6b7280', marginBottom: '1.5rem' }}>{error}</p>
          <button
            onClick={() => navigate('/login')}
            style={{
              padding: '0.5rem 1rem',
              backgroundColor: '#667eea',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer'
            }}
          >
            Back to Login
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '100vh',
      backgroundColor: '#f3f4f6',
      gap: '1rem'
    }}>
      <div style={{
        fontSize: '1.125rem',
        color: '#4b5563',
        fontWeight: 500
      }}>
        Processing authentication...
      </div>
      <div style={{
        fontSize: '0.875rem',
        color: '#6b7280'
      }}>
        Please wait while we complete your sign-in
      </div>
    </div>
  );
}
