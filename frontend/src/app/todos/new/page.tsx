// src/app/todos/new/page.tsx
'use client';

import { useRouter } from 'next/navigation';
import { useAuth } from '../../../hooks/use-auth';
import { useTodos } from '../../../hooks/use-todos';
import TodoForm from '@/components/todos/TodoForm';
import { CreateTodoRequest } from '../../../core/domain/todo';
import { useEffect } from 'react';

export default function NewTodoPage() {
  const { user, loading: authLoading } = useAuth();
  const { createTodo } = useTodos();
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    }
  }, [user, authLoading, router]);

  const handleSubmit = async (data: CreateTodoRequest) => {
    await createTodo(data);
    router.push('/');
  };

  if (authLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (!user) {
    return null; // Will redirect to login
  }

  return (
    <div>
      <h1 className="text-2xl font-semibold text-gray-900 mb-6">Add New Todo</h1>
      <TodoForm onSubmit={handleSubmit} />
    </div>
  );
}