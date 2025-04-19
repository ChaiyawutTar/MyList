'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useAuth } from '@/hooks/use-auth';
import { useTodos } from '@/hooks/use-todos';
import TodoForm from '@/components/todos/TodoForm';
import { Todo, UpdateTodoRequest } from '@/core/domain/todo';
import { Loading } from '@/components/ui/loading';
import { Card, CardContent } from '@/components/ui/card';
import { AlertCircle } from 'lucide-react';

export default function EditTodoPage() {
  const params = useParams();
  const { user, loading: authLoading } = useAuth();
  const { getTodoById, updateTodo } = useTodos();
  const [todo, setTodo] = useState<Todo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [imagePreloaded, setImagePreloaded] = useState(false);
  const router = useRouter();
  const todoId = params?.id ? parseInt(params.id as string) : NaN;

  // Fetch todo data
  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
      return;
    }
  
    const fetchTodo = async () => {
      if (!user || isNaN(todoId)) return;
      
      try {
        setLoading(true);
        const fetchedTodo = await getTodoById(todoId);
        setTodo(fetchedTodo);
      } catch (err) {
        setError('Failed to fetch todo');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
  
    if (user && !isNaN(todoId)) {
      fetchTodo();
    }
  }, [user, authLoading, todoId, router, getTodoById]);

  // Preload the image if it exists
  useEffect(() => {
    if (todo?.image_id && !imagePreloaded) {
      const API_URL = process.env.NEXT_PUBLIC_API_URL;
      const img = new Image();
      img.src = `${API_URL}/images/${todo.image_id}`;
      img.onload = () => setImagePreloaded(true);
    }
  }, [todo, imagePreloaded]);

  const handleSubmit = async (data: UpdateTodoRequest) => {
    try {
      setLoading(true);
      await updateTodo(todoId, data);
      router.push('/');
    } catch (err) {
      setError('Failed to update todo');
      console.error(err);
      setLoading(false);
    }
  };

  // Show loading state
  if (authLoading || (loading && !todo)) {
    return (
      <div className="flex justify-center items-center min-h-[60vh]">
        <Loading />
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <Card className="border-destructive">
        <CardContent className="pt-6">
          <div className="flex items-center gap-2 text-destructive">
            <AlertCircle className="h-5 w-5" />
            <p>{error}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Redirect if no user or todo
  if (!user || !todo) {
    return null;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold tracking-tight">Edit Todo</h1>
      <TodoForm
        initialData={{
          title: todo.title,
          description: todo.description,
          status: todo.status,
          image_id: todo.image_id,
        }}
        onSubmit={handleSubmit}
        isEditing
        isSubmitting={loading}
      />
    </div>
  );
}