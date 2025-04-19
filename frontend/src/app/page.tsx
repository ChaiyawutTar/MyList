// src/app/page.tsx
'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../hooks/use-auth';
import { Loading } from '@/components/ui/loading';

export default function Home() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      if (user) {
        // If user is logged in, redirect to /todos
        router.push('/todos');
      } else {
        // If user is not logged in, redirect to login
        router.push('/login');
      }
    }
  }, [user, loading, router]);

  if (loading) {
    return <Loading />;
  }

  return null; // Will redirect based on auth state
}