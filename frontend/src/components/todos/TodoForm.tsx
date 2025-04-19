'use client';

import { useState, FormEvent } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Loader2, X } from "lucide-react";
import type { CreateTodoRequest, UpdateTodoRequest, Todo } from '@/core/domain/todo';
import Image from 'next/image';
import { useEffect } from 'react';

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export interface TodoFormProps {
  initialData?: Partial<Todo>;
  onSubmit: (data: CreateTodoRequest | UpdateTodoRequest) => Promise<void>;
  isEditing?: boolean;
  isSubmitting?: boolean;
}

export default function TodoForm({ 
  initialData, 
  onSubmit, 
  isEditing = false,
  isSubmitting = false 
}: TodoFormProps) {
  const [title, setTitle] = useState(initialData?.title || '');
  const [description, setDescription] = useState(initialData?.description || '');
  const [status, setStatus] = useState<'pending' | 'in_progress' | 'done'>(initialData?.status || 'pending');
  const [image, setImage] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [imagePreview, setImagePreview] = useState<string | null>(null);

  // Generate preview when a new image is selected
  useEffect(() => {
    if (image) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
      };
      reader.readAsDataURL(image);
    } else {
      setImagePreview(null);
    }
  }, [image]);

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
        setImagePreview(null);
      }
    } catch (error) {
      console.error('Error submitting todo:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const clearImage = () => {
    setImage(null);
    setImagePreview(null);
  };

  const submitting = isLoading || isSubmitting;

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
              className="mb-2"
            />
            
            {/* Show image preview when a new image is selected */}
            {imagePreview && (
              <div className="mt-4">
                <div className="flex items-center justify-between mb-2">
                  <p className="text-sm text-muted-foreground">New image preview:</p>
                  <Button 
                    type="button" 
                    variant="outline" 
                    size="sm" 
                    onClick={clearImage}
                    className="h-8 px-2"
                  >
                    <X className="h-4 w-4 mr-1" />
                    Clear
                  </Button>
                </div>
                <div className="relative rounded-md overflow-hidden border">
                  <Image
                    src={imagePreview}
                    alt="Image preview"
                    width={300}
                    height={200}
                    className="object-contain w-full"
                    style={{ maxHeight: '200px' }}
                  />
                </div>
              </div>
            )}
            
            {/* Show current image if editing and no new image is selected */}
            {!imagePreview && initialData?.image_id && (
              <div className="mt-4">
                <p className="text-sm text-muted-foreground mb-2">Current image:</p>
                <div className="relative rounded-md overflow-hidden border">
                  <Image
                    src={`${API_URL}/images/${initialData.image_id}`}
                    alt={`Image for ${initialData.title}`}
                    width={300}
                    height={200}
                    className="object-contain w-full"
                    style={{ maxHeight: '200px' }}
                    loading="lazy"
                  />
                </div>
              </div>
            )}
          </div>
        </CardContent>
        <CardFooter>
          <Button type="submit" disabled={submitting}>
            {submitting ? (
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