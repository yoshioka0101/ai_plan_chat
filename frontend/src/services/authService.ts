import type { User } from '../types/auth';

export const authService = {
  // Logout (client-side)
  logout(): void {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  // Token management
  setToken(token: string): void {
    localStorage.setItem('token', token);
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  },

  // User management
  setUser(user: User): void {
    localStorage.setItem('user', JSON.stringify(user));
  },

  getUser(): User | null {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },
};
