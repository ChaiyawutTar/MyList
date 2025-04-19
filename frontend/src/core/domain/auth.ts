
export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
  oauth_provider?: string;
  picture?: string;
}

// Add OAuth login method
export interface OAuthLoginParams {
  provider: string;
}
  
  export interface LoginRequest {
    email: string;
    password: string;
  }
  
  export interface SignupRequest {
    username: string;
    email: string;
    password: string;
  }
  
  export interface AuthResponse {
    token: string;
    user: User;
  }