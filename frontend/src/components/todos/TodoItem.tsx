// src/components/todos/TodoItem.tsx
'use client';

import { useState } from 'react';
import { Todo } from '@/core/domain/todo';
import Link from 'next/link';

interface TodoItemProps {
  todo: Todo;
  onDelete: (id: number) => Promise<void>;
}

export default function TodoItem({ todo, onDelete }: TodoItemProps) {
  const [isLoading, setIsLoading] = useState(false);

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this todo?')) {
      setIsLoading(true);
      try {
        await onDelete(todo.id);
      } catch (error) {
        console.error('Error deleting todo:', error);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const getStatusColor = () => {
    switch (todo.status) {
      case 'done':
        return 'bg-green-100 text-green-800';
      case 'in_progress':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="bg-white shadow rounded-lg p-4 mb-4">
      <div className="flex justify-between">
        <div>
          <h3 className="text-lg font-medium text-gray-900">{todo.title}</h3>
          <p className="mt-1 text-sm text-gray-500">{todo.description}</p>
          <span className={`inline-flex mt-2 items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor()}`}>
            {todo.status.replace('_', ' ')}
          </span>
        </div>
        <div>
          <Link
            href={`/todos/${todo.id}/edit`}
            className="mr-2 text-indigo-600 hover:text-indigo-900"
          >
            Edit
          </Link>
          <button
            onClick={handleDelete}
            className="text-red-600 hover:text-red-900"
            disabled={isLoading}
          >
            Delete
          </button>
        </div>
      </div>
      {todo.image_path && (
        <div className="mt-4">
          <img
            src={`http://localhost:8080/${todo.image_path}`}
            alt={todo.title}
            className="h-40 w-auto object-cover rounded"
          />
        </div>
      )}
    </div>
  );
}