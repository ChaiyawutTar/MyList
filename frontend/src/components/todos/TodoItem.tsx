// src/components/todos/TodoItem.tsx
'use client';

import { useState, useEffect, useRef } from 'react';
import { Todo } from '@/core/domain/todo';
import Link from 'next/link';
import Image from 'next/image';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Edit, Trash2, Loader2, AlertTriangle, ImageIcon } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface TodoItemProps {
  todo: Todo;
  onDelete: (id: number) => Promise<void>;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function TodoItem({ todo, onDelete }: TodoItemProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);
  const [imageError, setImageError] = useState(false);
  const [shouldLoadImage, setShouldLoadImage] = useState(false);
  const todoRef = useRef<HTMLDivElement>(null);

  // Set up intersection observer for lazy loading
  useEffect(() => {
    if (!todo.image_id) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          setShouldLoadImage(true);
          observer.disconnect();
        }
      },
      { 
        rootMargin: '200px', // Start loading when within 200px of viewport
        threshold: 0.1 
      }
    );

    if (todoRef.current) {
      observer.observe(todoRef.current);
    }

    return () => {
      observer.disconnect();
    };
  }, [todo.image_id]);

  const handleDelete = async () => {
    setIsLoading(true);
    try {
      await onDelete(todo.id);
      setShowDeleteDialog(false);
    } catch (error) {
      console.error('Error deleting todo:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const getStatusColor = () => {
    switch (todo.status) {
      case 'done':
        return 'bg-green-100 text-green-800 hover:bg-green-200';
      case 'in_progress':
        return 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200';
      default:
        return 'bg-gray-100 text-gray-800 hover:bg-gray-200';
    }
  };

  return (
    <>
      <Card ref={todoRef}>
        <CardHeader className="pb-2">
          <div className="flex justify-between items-start">
            <CardTitle className="text-lg">{todo.title}</CardTitle>
            <Badge variant="outline" className={getStatusColor()}>
              {todo.status.replace('_', ' ')}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground line-clamp-3">
            {todo.description || "No description provided"}
          </p>
          
          {todo.image_id && (
            <div className="mt-4 relative h-40 bg-gray-100 rounded-md overflow-hidden">
              {!shouldLoadImage ? (
                <div className="absolute inset-0 flex items-center justify-center">
                  <ImageIcon className="h-8 w-8 text-gray-400" />
                </div>
              ) : !imageLoaded && !imageError ? (
                <div className="absolute inset-0 flex items-center justify-center">
                  <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                </div>
              ) : null}
              
              {shouldLoadImage && (
                <Image
                  width={400}
                  height={200}
                  src={`${API_URL}/images/${todo.image_id}`}
                  alt={todo.title}
                  className={`h-40 w-full object-cover transition-opacity duration-300 ${imageLoaded ? 'opacity-100' : 'opacity-0'}`}
                  onLoad={() => setImageLoaded(true)}
                  onError={() => {
                    setImageError(true);
                    console.error(`Failed to load image for todo ${todo.id}`);
                  }}
                  priority={false}
                  loading="lazy"
                />
              )}
              
              {imageError && (
                <div className="absolute inset-0 flex flex-col items-center justify-center">
                  <AlertTriangle className="h-6 w-6 text-red-500 mb-2" />
                  <span className="text-sm text-red-500">Failed to load image</span>
                </div>
              )}
            </div>
          )}
        </CardContent>
        <CardFooter className="flex justify-between pt-2">
          <Button variant="outline" size="sm" asChild>
            <Link href={`/todos/${todo.id}/edit`}>
              <Edit className="mr-2 h-4 w-4" />
              Edit
            </Link>
          </Button>
          <Button 
            variant="destructive" 
            size="sm" 
            onClick={() => setShowDeleteDialog(true)}
          >
            <Trash2 className="mr-2 h-4 w-4" />
            Delete
          </Button>
        </CardFooter>
      </Card>

      {/* Delete Confirmation Dialog */}
      <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-destructive" />
              Confirm Deletion
            </DialogTitle>
            <DialogDescription>
              Are you sure you want to delete &quot;{todo.title}&quot;? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0">
            <Button
              variant="outline"
              onClick={() => setShowDeleteDialog(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={isLoading}
            >
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Deleting...
                </>
              ) : (
                'Delete'
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}