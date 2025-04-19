// src/components/todos/TodoForm.tsx
'use client';

import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Loader2 } from "lucide-react";
// Remove the Image import if using regular img tag

export default function TodoForm({ initialData, onSubmit, isEditing = false }) {
  const [title, setTitle] = useState(initialData?.title || '');
  const [description, setDescription] = useState(initialData?.description || '');
  const [status, setStatus] = useState(initialData?.status || 'pending');
  const [image, setImage] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
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
        // Reset form if creating new todo
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
            <Select 
              value={status} 
              onValueChange={setStatus}
            >
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
            {initialData?.image_path && (
              <div className="mt-2">
                <p className="text-sm text-muted-foreground">Current image:</p>
                {/* Using regular img tag instead of Next.js Image */}
                <img
                  src={`http://localhost:8080/${initialData.image_path}`}
                  alt={`Image for ${initialData.title}`}
                  className="mt-2 h-20 w-auto object-cover rounded"
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