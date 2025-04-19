import { AuthResponse, LoginRequest, SignupRequest, User } from '@/core/domain/auth';

export interface AuthRepository {
  login(request: LoginRequest): Promise<AuthResponse>;
  signup(request: SignupRequest): Promise<AuthResponse>;
  oauthLogin(provider: string): void; // Redirects to OAuth provider
  getCurrentUser(): User | null;
  saveToken(token: string): void;
  getToken(): string | null;
  removeToken(): void;
  isAuthenticated(): boolean;
}