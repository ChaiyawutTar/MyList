// // src/app/(auth)/callback/oauth-callback-client.tsx
// 'use client';

// import { useEffect, useState } from 'react';
// import { useRouter, useSearchParams } from 'next/navigation';
// import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
// import { Loader2 } from "lucide-react";

// export default function OAuthCallbackClient() {
//     const router = useRouter();
//     const searchParams = useSearchParams();
//     const token = searchParams.get('token');
//     const [message, setMessage] = useState('Processing authentication...');
//     const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');

//     useEffect(() => {
//         if (!token) {
//             setMessage('Authentication failed. No token received.');
//             setStatus('error');
//             setTimeout(() => router.push('/login?error=oauth_failed'), 2000);
//             return;
//         }

//         try {
//             // Store the token
//             localStorage.setItem('token', token);

//             // Update message
//             setMessage('Authentication successful! Redirecting to your todos...');
//             setStatus('success');

//             // Redirect to todos page after a short delay
//             setTimeout(() => router.push('/todos'), 1000);
//         } catch (error) {
//             console.error('Error storing token:', error);
//             setMessage('Authentication failed. Please try again.');
//             setStatus('error');
//             setTimeout(() => router.push('/login?error=storage_failed'), 2000);
//         }
//     }, [token, router]);

//     return (
//         <div className="flex flex-col items-center justify-center min-h-screen p-4">
//             <Card className="w-full max-w-md">
//                 <CardHeader>
//                     <CardTitle>OAuth Authentication</CardTitle>
//                     <CardDescription>
//                         {status === 'loading' && 'Completing your authentication...'}
//                         {status === 'success' && 'Authentication completed successfully!'}
//                         {status === 'error' && 'There was a problem with authentication.'}
//                     </CardDescription>
//                 </CardHeader>
//                 <CardContent className="flex flex-col items-center justify-center py-6">
//                     {status === 'loading' && <Loader2 className="h-8 w-8 animate-spin text-primary" />}
//                     <p className="mt-4 text-center text-muted-foreground">{message}</p>
//                 </CardContent>
//             </Card>
//         </div>
//     );
// }