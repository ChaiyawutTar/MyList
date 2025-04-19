// src/components/todos/TodoList.tsx
'use client';

import { useEffect, useRef, useState } from 'react';
import { useTodos } from '../../hooks/use-todos';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { PlusCircle, Loader2 } from "lucide-react";
import Link from 'next/link';
import TodoItem from './TodoItem';
import { useInView } from 'react-intersection-observer';

export default function TodoList() {
  const { todos, loading, error, fetchTodos, deleteTodo } = useTodos();
  const [visibleCount, setVisibleCount] = useState(9); // Start with 9 items (3x3 grid)
  const containerRef = useRef<HTMLDivElement>(null);
  
  // Set up intersection observer for infinite scrolling
  const { ref: loadMoreRef, inView } = useInView({
    threshold: 0.1,
    rootMargin: '200px',
  });

  // Load more items when the user scrolls to the bottom
  useEffect(() => {
    if (inView && todos && visibleCount < todos.length) {
      setVisibleCount(prev => Math.min(prev + 6, todos.length));
    }
  }, [inView, todos, visibleCount]);

  useEffect(() => {
    fetchTodos();
  }, [fetchTodos]);

  // Reset visible count when todos change
  useEffect(() => {
    if (todos) {
      setVisibleCount(Math.min(9, todos.length));
    }
  }, [todos]);

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
        <div ref={containerRef} className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {todos.slice(0, visibleCount).map((todo) => (
              <TodoItem key={todo.id} todo={todo} onDelete={deleteTodo} />
            ))}
          </div>
          
          {/* Load more trigger element */}
          {todos.length > visibleCount && (
            <div 
              ref={loadMoreRef} 
              className="flex justify-center py-4"
            >
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
            </div>
          )}
          
          {/* Show total count */}
          <div className="text-center text-sm text-muted-foreground">
            Showing {visibleCount} of {todos.length} todos
          </div>
        </div>
      )}
    </div>
  );
}