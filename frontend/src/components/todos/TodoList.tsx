// src/components/todos/TodoList.tsx
'use client';

import { useEffect } from 'react';
import TodoItem from './TodoItem';
import { useTodos } from '../../hooks/use-todos';
import Link from 'next/link';

export default function TodoList() {
  const { todos, loading, error, fetchTodos, deleteTodo } = useTodos();

  useEffect(() => {
    fetchTodos();
  }, [fetchTodos]);

  if (loading && todos.length === 0) {
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

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold text-gray-900">My Todos</h1>
        <Link
          href="/todos/new"
          className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
        >
          Add New Todo
        </Link>
      </div>

      {todos.length === 0 ? (
        <div className="bg-white shadow rounded-lg p-6 text-center">
          <p className="text-gray-500">No todos yet. Add your first todo!</p>
        </div>
      ) : (
        todos.map((todo) => (
          <TodoItem key={todo.id} todo={todo} onDelete={deleteTodo} />
        ))
      )}
    </div>
  );
}