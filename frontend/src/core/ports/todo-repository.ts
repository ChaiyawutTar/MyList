// src/core/ports/todo-repository.ts
import { CreateTodoRequest, Todo, UpdateTodoRequest } from '@/core/domain/todo';

export interface TodoRepository {
  getAllTodos(): Promise<Todo[]>;
  getTodoById(id: number): Promise<Todo>;
  createTodo(request: CreateTodoRequest): Promise<Todo>;
  updateTodo(id: number, request: UpdateTodoRequest): Promise<Todo>;
  deleteTodo(id: number): Promise<void>;
}