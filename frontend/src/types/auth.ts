export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  email: string;
  password: string;
  name: string;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  logout: () => void;
  setAuth: (user: User, token: string) => void;
  isAuthenticated: boolean;
  isLoading: boolean;
}
