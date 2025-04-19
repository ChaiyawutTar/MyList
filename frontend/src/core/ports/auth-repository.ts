import { AuthResponse, LoginRequest, SignupRequest, User } from '@/core/domain/auth';

export interface AuthRepository {
  login(request: LoginRequest): Promise<AuthResponse>;
  signup(request: SignupRequest): Promise<AuthResponse>;
  oauthLogin(provider: string, code?: string): Promise<AuthResponse>;
  getCurrentUser(): User | null;
  saveToken(token: string): void;
  getToken(): string | null;
  removeToken(): void;
  isAuthenticated(): boolean;
}