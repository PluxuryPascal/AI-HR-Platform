"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { useState, type ReactNode } from "react";

function makeQueryClient() {
    return new QueryClient({
        defaultOptions: {
            queries: {
                // Data considered fresh for 60 seconds — avoids unnecessary refetches
                staleTime: 60 * 1000,
                // Retry failed requests once before surfacing the error
                retry: 1,
                // Don't refetch every time the browser tab regains focus
                refetchOnWindowFocus: false,
            },
        },
    });
}

let browserQueryClient: QueryClient | undefined;

function getQueryClient() {
    // Server: always create a new client (avoids sharing state across requests)
    if (typeof window === "undefined") return makeQueryClient();

    // Browser: reuse a single client instance
    if (!browserQueryClient) browserQueryClient = makeQueryClient();
    return browserQueryClient;
}

export function QueryProvider({ children }: { children: ReactNode }) {
    const [queryClient] = useState(getQueryClient);

    return (
        <QueryClientProvider client={queryClient}>
            {children}
            <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
    );
}
