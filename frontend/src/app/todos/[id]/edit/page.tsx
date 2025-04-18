'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useAuth } from '@/hooks/use-auth';
import { useTodos } from '@/hooks/use-todos';
import TodoForm from '@/components/todos/TodoForm';
import { Todo, UpdateTodoRequest } from '@/core/domain/todo';

export default function EditTodoPage() {
  const params = useParams();
  const { user, loading: authLoading } = useAuth();
  const { getTodoById, updateTodo } = useTodos();
  const [todo, setTodo] = useState<Todo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const todoId = params?.id ? parseInt(params.id as string) : NaN;

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
      return;
    }

    if (user && !isNaN(todoId)) {
      fetchTodo();
    }
  }, [user, authLoading, todoId, router]);

  const fetchTodo = async () => {
    try {
      const fetchedTodo = await getTodoById(todoId);
      setTodo(fetchedTodo);
    } catch (err) {
      setError('Failed to fetch todo');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (data: UpdateTodoRequest) => {
    await updateTodo(todoId, data);
    router.push('/');
  };

  if (authLoading || loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
        {error}
      </div>
    );
  }

  if (!user || !todo) {
    return null; // Will redirect to login or home
  }

  return (
    <div>
      <h1 className="text-2xl font-semibold text-gray-900 mb-6">Edit Todo</h1>
      <TodoForm
        initialData={{
          title: todo.title,
          description: todo.description,
          status: todo.status,
          image_path: todo.image_path,
        }}
        onSubmit={handleSubmit}
        isEditing
      />
    </div>
  );
}