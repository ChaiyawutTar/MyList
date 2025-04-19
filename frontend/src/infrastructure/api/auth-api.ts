import { AuthRepository } from '@/core/ports/auth-repository';
import { AuthResponse, LoginRequest, SignupRequest, User } from '@/core/domain/auth';
import { apiClient } from '@/infrastructure/api/api-client';
import { jwtDecode } from 'jwt-decode';

export class AuthRepositoryImpl implements AuthRepository {
  async login(request: LoginRequest): Promise<AuthResponse> {
    // Convert LoginRequest to Record<string, unknown> to satisfy the type constraint
    return apiClient.post<AuthResponse>('/login', request as unknown as Record<string, unknown>);
  }

  async signup(request: SignupRequest): Promise<AuthResponse> {
    // Convert SignupRequest to Record<string, unknown> to satisfy the type constraint
    return apiClient.post<AuthResponse>('/signup', request as unknown as Record<string, unknown>);
  }

  // Implement the missing oauthLogin method
  async oauthLogin(provider: string, code: string): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/oauth/callback', {
      provider,
      code
    } as Record<string, unknown>);
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
    } catch (err) {
      // Use err instead of error to avoid the unused variable warning
      console.error('Error decoding token:', err);
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
    } catch (err) {
      // Use err instead of error to avoid the unused variable warning
      console.error('Error checking authentication:', err);
      this.removeToken();
      return false;
    }
  }
}

export const authRepository = new AuthRepositoryImpl();