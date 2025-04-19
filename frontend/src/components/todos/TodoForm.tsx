'use client';

import { useState, FormEvent } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Loader2 } from "lucide-react";
import type { CreateTodoRequest, UpdateTodoRequest, Todo } from '@/core/domain/todo';
import Image from 'next/image';

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export interface TodoFormProps {
  initialData?: Partial<Todo>;
  onSubmit: (data: CreateTodoRequest | UpdateTodoRequest) => Promise<void>;
  isEditing?: boolean;
}

export default function TodoForm({ initialData, onSubmit, isEditing = false }: TodoFormProps) {
  const [title, setTitle] = useState(initialData?.title || '');
  const [description, setDescription] = useState(initialData?.description || '');
  const [status, setStatus] = useState<'pending' | 'in_progress' | 'done'>(initialData?.status || 'pending');
  const [image, setImage] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await onSubmit({
        title,
        description,
        status,
        image: image || undefined,
      });

      if (!isEditing) {
        setTitle('');
        setDescription('');
        setStatus('pending');
        setImage(null);
      }
    } catch (error) {
      console.error('Error submitting todo:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{isEditing ? 'Edit Todo' : 'Add New Todo'}</CardTitle>
        <CardDescription>
          {isEditing 
            ? 'Update your todo item details' 
            : 'Fill in the details to create a new todo item'}
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter todo title"
              required
            />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Enter todo description"
              rows={3}
            />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <Select value={status} onValueChange={(val: 'pending' | 'in_progress' | 'done') => setStatus(val)}>
              <SelectTrigger>
                <SelectValue placeholder="Select status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="in_progress">In Progress</SelectItem>
                <SelectItem value="done">Done</SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="image">Image (Optional)</Label>
            <Input
              id="image"
              type="file"
              onChange={(e) => setImage(e.target.files?.[0] || null)}
              accept="image/*"
            />
            {initialData?.image_id && (
              <div className="mt-2">
                <p className="text-sm text-muted-foreground">Current image:</p>
                <Image
                  src={`${API_URL}/images/${initialData.image_id}`}
                  alt={`Image for ${initialData.title}`}
                  width={80} // adjust based on expected size
                  height={80}
                  className="mt-2 object-cover rounded"
                  style={{ height: 'auto' }} // optional
                />
              </div>
            )}
          </div>
        </CardContent>
        <CardFooter>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                {isEditing ? 'Saving...' : 'Adding...'}
              </>
            ) : (
              isEditing ? 'Save Changes' : 'Add Todo'
            )}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}
