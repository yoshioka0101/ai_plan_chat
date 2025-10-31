import type { User } from '../types/auth';

// localStorageが利用可能かチェック（SSR対応）
const isLocalStorageAvailable = (): boolean => {
  return typeof window !== 'undefined' && typeof window.localStorage !== 'undefined';
};

export const authService = {
  // Logout (client-side)
  logout(): void {
    if (!isLocalStorageAvailable()) return;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  // Token management
  setToken(token: string): void {
    if (!isLocalStorageAvailable()) return;
    localStorage.setItem('token', token);
  },

  getToken(): string | null {
    if (!isLocalStorageAvailable()) return null;
    return localStorage.getItem('token');
  },

  // User management
  setUser(user: User): void {
    if (!isLocalStorageAvailable()) return;
    try {
      localStorage.setItem('user', JSON.stringify(user));
    } catch (error) {
      console.error('Failed to save user to localStorage:', error);
    }
  },

  getUser(): User | null {
    if (!isLocalStorageAvailable()) return null;
    
    const userStr = localStorage.getItem('user');
    if (!userStr) return null;

    try {
      return JSON.parse(userStr) as User;
    } catch (error) {
      console.error('Failed to parse user from localStorage:', error);
      // 不正なJSONの場合は削除
      localStorage.removeItem('user');
      return null;
    }
  },
};
