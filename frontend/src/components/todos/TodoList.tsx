// src/components/todos/TodoList.tsx
'use client';

import { useEffect } from 'react';
// import { useRouter } from 'next/navigation';
import { useTodos } from '../../hooks/use-todos';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { PlusCircle, Loader2 } from "lucide-react";
import Link from 'next/link';
import TodoItem from './TodoItem';

export default function TodoList() {
  const { todos, loading, error, fetchTodos, deleteTodo } = useTodos();
  // const router = useRouter();

  useEffect(() => {
    fetchTodos();
  }, [fetchTodos]);

  if (loading && (!todos || todos.length === 0)) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight">My Todos</h1>
        <Button asChild>
          <Link href="/todos/new">
            <PlusCircle className="mr-2 h-4 w-4" />
            Add New Todo
          </Link>
        </Button>
      </div>

      {error && (
        <Card className="border-destructive">
          <CardContent className="pt-6">
            <p className="text-destructive">{error}</p>
          </CardContent>
        </Card>
      )}

      {!todos || todos.length === 0 ? (
        <Card>
          <CardHeader>
            <CardTitle>No todos yet</CardTitle>
            <CardDescription>
              Get started by adding your first todo item
            </CardDescription>
          </CardHeader>
          <CardContent className="flex justify-center pb-6">
            <Button asChild>
              <Link href="/todos/new">
                <PlusCircle className="mr-2 h-4 w-4" />
                Create Todo
              </Link>
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {todos.map((todo) => (
            <TodoItem key={todo.id} todo={todo} onDelete={deleteTodo} />
          ))}
        </div>
      )}
    </div>
  );
}