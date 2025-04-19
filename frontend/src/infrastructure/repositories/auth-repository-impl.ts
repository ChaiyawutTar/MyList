// src/infrastructure/repositories/auth-repository-impl.ts
import { AuthRepository } from '@/core/ports/auth-repository';
import { AuthResponse, LoginRequest, SignupRequest, User } from '@/core/domain/auth';
import { apiClient } from '@/infrastructure/api/api-client';
import { jwtDecode } from 'jwt-decode';

export class AuthRepositoryImpl implements AuthRepository {
  async login(request: LoginRequest): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/login', request);
  }

  async signup(request: SignupRequest): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/signup', request);
  }

  oauthLogin(provider: string): void {
    if (typeof window !== 'undefined') {
        window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/${provider}`;
    }
}

  getCurrentUser(): User | null {
    const token = this.getToken();
    if (!token) return null;

    try {
      const decoded = jwtDecode<{ user_id: number }>(token);
      // In a real app, you might want to fetch the user details from the API
      // For simplicity, we're just returning the user ID
      return {
        id: decoded.user_id,
        username: '',
        email: '',
        created_at: '',
      };
    } catch (error) {
      this.removeToken();
      return null;
    }
  }

  saveToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }

  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('token');
    }
    return null;
  }

  removeToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }

  isAuthenticated(): boolean {
    const token = this.getToken();
    if (!token) return false;

    try {
      const decoded = jwtDecode<{ exp: number }>(token);
      return decoded.exp * 1000 > Date.now();
    } catch (error) {
      this.removeToken();
      return false;
    }
  }
}

export const authRepository = new AuthRepositoryImpl();