'use client';

import { useState } from 'react';

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
    <div className="bg-white shadow rounded-lg p-6 mb-6">
      <h2 className="text-lg font-medium text-gray-900 mb-4">
        {isEditing ? 'Edit Todo' : 'Add New Todo'}
      </h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="title" className="block text-gray-700 text-sm font-bold mb-2">Title</label>
          <input
            id="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Enter todo title"
            required
            aria-required="true"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="description" className="block text-gray-700 text-sm font-bold mb-2">Description</label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Enter todo description"
            rows={3}
          />
        </div>
        <div className="mb-4">
          <label htmlFor="status" className="block text-gray-700 text-sm font-bold mb-2">Status</label>
          <select
            id="status"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            aria-label="Todo status"
          >
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="done">Done</option>
          </select>
        </div>
        <div className="mb-4">
          <label htmlFor="image" className="block text-gray-700 text-sm font-bold mb-2">Image (Optional)</label>
          <input
            id="image"
            type="file"
            onChange={(e) => setImage(e.target.files?.[0] || null)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            accept="image/*"
            aria-label="Upload todo image"
          />
          {initialData?.image_path && (
            <div className="mt-2">
              <p className="text-sm text-gray-500">Current image:</p>
              <img
                src={`http://localhost:8080/${initialData.image_path}`}
                alt={`Image for ${initialData.title}`}
                className="mt-2 h-20 w-auto object-cover rounded"
              />
            </div>
          )}
        </div>
        <div className="flex justify-end">
          <button
            type="submit"
            className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
            disabled={isLoading}
          >
            {isLoading ? (isEditing ? 'Saving...' : 'Adding...') : (isEditing ? 'Save' : 'Add Todo')}
          </button>
        </div>
      </form>
    </div>
  );
}
