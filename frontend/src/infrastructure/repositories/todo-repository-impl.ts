// src/infrastructure/repositories/todo-repository-impl.ts
import { TodoRepository } from '@/core/ports/todo-repository';
import { CreateTodoRequest, Todo, UpdateTodoRequest } from '@/core/domain/todo';
import { apiClient } from '@/infrastructure/api/api-client';

export class TodoRepositoryImpl implements TodoRepository {
  async getAllTodos(): Promise<Todo[]> {
    return apiClient.get<Todo[]>('/todos');
  }
  
  async getTodoById(id: number): Promise<Todo> {
    return apiClient.get<Todo>(`/todos/${id}`);
  }

  async createTodo(request: CreateTodoRequest): Promise<Todo> {
    const formData = new FormData();
    formData.append('title', request.title);
    formData.append('description', request.description);
    formData.append('status', request.status);
    if (request.image) {
      formData.append('image', request.image);
    }

    return apiClient.postFormData<Todo>('/todos', formData);
  }

  async updateTodo(id: number, request: UpdateTodoRequest): Promise<Todo> {
    const formData = new FormData();
    formData.append('title', request.title);
    formData.append('description', request.description);
    formData.append('status', request.status);
    if (request.image) {
      formData.append('image', request.image);
    }

    return apiClient.putFormData<Todo>(`/todos/${id}`, formData);
  }

  async deleteTodo(id: number): Promise<void> {
    return apiClient.delete<void>(`/todos/${id}`);
  }
}

export const todoRepository = new TodoRepositoryImpl();