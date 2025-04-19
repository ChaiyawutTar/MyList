
import { AuthRepository } from '@/core/ports/auth-repository';
import { AuthResponse, LoginRequest, SignupRequest, User } from '@/core/domain/auth';

export class AuthUseCases {
  constructor(private authRepository: AuthRepository) {}

  async login(request: LoginRequest): Promise<AuthResponse> {
    const response = await this.authRepository.login(request);
    this.authRepository.saveToken(response.token);
    return response;
  }

  async signup(request: SignupRequest): Promise<AuthResponse> {
    const response = await this.authRepository.signup(request);
    this.authRepository.saveToken(response.token);
    return response;
  }

  getCurrentUser(): User | null {
    return this.authRepository.getCurrentUser();
  }

  isAuthenticated(): boolean {
    return this.authRepository.isAuthenticated();
  }

  logout(): void {
    this.authRepository.removeToken();
  }

  oauthLogin(provider: string): void {
    this.authRepository.oauthLogin(provider);
  }
}