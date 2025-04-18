
export interface Todo {
    id: number;
    user_id: number;
    title: string;
    description: string;
    status: 'pending' | 'in_progress' | 'done';
    image_path?: string;
    created_at: string;
    updated_at: string;
  }
  
  export interface CreateTodoRequest {
    title: string;
    description: string;
    status: string;
    image?: File;
  }
  
  export interface UpdateTodoRequest {
    title: string;
    description: string;
    status: string;
    image?: File;
  }