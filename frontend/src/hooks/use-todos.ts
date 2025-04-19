// src/hooks/use-todos.ts
'use client';

import { useState, useCallback, useMemo } from 'react';
import { Todo, CreateTodoRequest, UpdateTodoRequest } from '@/core/domain/todo';
import { TodoUseCases } from '@/core/usecases/todo-usecases';
import { todoRepository } from '@/infrastructure/repositories/todo-repository-impl';

export const useTodos = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Use useMemo to create todoUseCases once and avoid recreating it on every render
  const todoUseCases = useMemo(() => new TodoUseCases(todoRepository), []);

  const fetchTodos = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const fetchedTodos = await todoUseCases.getAllTodos();
      setTodos(fetchedTodos);
    } catch (err) {
      setError('Failed to fetch todos');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [todoUseCases]); // Add todoUseCases to the dependency array

  const getTodoById = useCallback(async (id: number) => {
    setLoading(true);
    setError(null);
    try {
      // Option 1: Use the existing todos array if it's already loaded
      if (todos.length > 0) {
        const todo = todos.find(t => t.id === id);
        if (todo) {
          return todo;
        }
      }
      
      // Option 2: Fetch all todos and find the one we need
      // This is a workaround until you implement the GET /todos/{id} endpoint
      const allTodos = await todoRepository.getAllTodos();
      const todo = allTodos.find(t => t.id === id);
      if (!todo) {
        throw new Error('Todo not found');
      }
      return todo;
    } catch (err) {
      setError('Failed to fetch todo');
      console.error(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [todos]); // This one is fine as it doesn't use todoUseCases

  const createTodo = useCallback(async (todo: CreateTodoRequest) => {
    setLoading(true);
    setError(null);
    try {
      const newTodo = await todoUseCases.createTodo(todo);
      setTodos((prevTodos) => [newTodo, ...prevTodos]);
      return newTodo;
    } catch (err) {
      setError('Failed to create todo');
      console.error(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [todoUseCases]); // Add todoUseCases to the dependency array

  const updateTodo = useCallback(async (id: number, todo: UpdateTodoRequest) => {
    setLoading(true);
    setError(null);
    try {
      const updatedTodo = await todoUseCases.updateTodo(id, todo);
      setTodos((prevTodos) =>
        prevTodos.map((t) => (t.id === id ? updatedTodo : t))
      );
      return updatedTodo;
    } catch (err) {
      setError('Failed to update todo');
      console.error(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [todoUseCases]); // Add todoUseCases to the dependency array

  const deleteTodo = useCallback(async (id: number) => {
    setLoading(true);
    setError(null);
    try {
      await todoUseCases.deleteTodo(id);
      setTodos((prevTodos) => prevTodos.filter((t) => t.id !== id));
    } catch (err) {
      setError('Failed to delete todo');
      console.error(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [todoUseCases]); // Add todoUseCases to the dependency array

  return {
    todos,
    loading,
    error,
    fetchTodos,
    getTodoById,
    createTodo,
    updateTodo,
    deleteTodo,
  };
};