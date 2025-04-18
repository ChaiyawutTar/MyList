// src/providers/auth-provider.tsx
'use client';

import { createContext, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { User } from '@/core/domain/auth';
import { AuthUseCases } from '@/core/usecases/auth-usecases';
import { authRepository } from '@/infrastructure/repositories/auth-repository-impl';

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

export const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  signup: async () => {},
  logout: () => {},
  loading: true,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const authUseCases = new AuthUseCases(authRepository);

  useEffect(() => {
    // Check if user is logged in
    const currentUser = authUseCases.getCurrentUser();
    setUser(currentUser);
    setLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    try {
      const response = await authUseCases.login({ email, password });
      setUser(response.user);
      router.push('/');
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  };

  const signup = async (username: string, email: string, password: string) => {
    try {
      const response = await authUseCases.signup({ username, email, password });
      setUser(response.user);
      router.push('/');
    } catch (error) {
      console.error('Signup error:', error);
      throw error;
    }
  };

  const logout = () => {
    authUseCases.logout();
    setUser(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, login, signup, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};