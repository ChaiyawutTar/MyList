// src/components/CallbackClient.tsx (adjust path as needed)
'use client'; // Essential: Marks this as a Client Component

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2 } from "lucide-react";

export default function CallbackClient() {
    const router = useRouter();
    const searchParams = useSearchParams(); // Now safe within a 'use client' component intended for Suspense
    const token = searchParams.get('token');
    const [message, setMessage] = useState('Processing authentication...');
    const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');

    useEffect(() => {
        if (!token) {
            setMessage('Authentication failed. No token received.');
            setStatus('error');
            setTimeout(() => router.push('/login?error=oauth_failed'), 2000);
            return;
        }

        try {
            // Store the token
            localStorage.setItem('token', token);

            // Update message
            setMessage('Authentication successful! Redirecting to your todos...'); // Updated message
            setStatus('success');

            // Redirect to /todos page after a short delay
            // ***** CHANGE MADE HERE *****
            setTimeout(() => router.push('/todos'), 1500); // Redirect to /todos

        } catch (error) {
            console.error('Error storing token:', error);
            // Handle potential localStorage errors (e.g., private Browse, storage full)
            setMessage('Authentication failed. Could not save session. Please try again.');
            setStatus('error');
            setTimeout(() => router.push('/login?error=storage_failed'), 2000);
        }
        // Only include dependencies that actually affect the effect's logic.
        // searchParams itself can change, but we only care about the 'token' derived from it.
        // Adding router ensures useEffect runs again if the router instance changes (unlikely but technically possible).
    }, [token, router]); // Dependency array updated

    return (
        <Card className="w-full max-w-md">
            <CardHeader>
                <CardTitle>OAuth Authentication</CardTitle>
                <CardDescription>
                    {status === 'loading' && 'Completing your authentication...'}
                    {status === 'success' && 'Authentication completed successfully!'}
                    {status === 'error' && 'There was a problem with authentication.'}
                </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col items-center justify-center py-6">
                {status === 'loading' && <Loader2 className="h-8 w-8 animate-spin text-primary" />}
                {/* Optionally show success/error icons */}
                {/* {status === 'success' && <CheckCircle2 className="h-8 w-8 text-green-500" />} */}
                {/* {status === 'error' && <XCircle className="h-8 w-8 text-red-500" />} */}
                <p className="mt-4 text-center text-muted-foreground">{message}</p>
            </CardContent>
        </Card>
    );
}