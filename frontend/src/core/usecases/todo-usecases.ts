// src/core/usecases/todo-usecases.ts
import { TodoRepository } from '../ports/todo-repository';
import { CreateTodoRequest, Todo, UpdateTodoRequest } from '@/core/domain/todo';

export class TodoUseCases {
  constructor(private todoRepository: TodoRepository) {}

  async getAllTodos(): Promise<Todo[]> {
    return this.todoRepository.getAllTodos();
  }
  
  async getTodoById(id: number): Promise<Todo> {
    return this.todoRepository.getTodoById(id);
  }

  async createTodo(request: CreateTodoRequest): Promise<Todo> {
    return this.todoRepository.createTodo(request);
  }

  async updateTodo(id: number, request: UpdateTodoRequest): Promise<Todo> {
    return this.todoRepository.updateTodo(id, request);
  }

  async deleteTodo(id: number): Promise<void> {
    return this.todoRepository.deleteTodo(id);
  }
}