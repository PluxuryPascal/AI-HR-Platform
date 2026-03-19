"use client";

import { useEffect, type ReactNode } from "react";
import { useGetProfile } from "@/features/profile/api/use-profile";
import { useAuth } from "@/store/use-auth";

export function AuthProvider({ children }: { children: ReactNode }) {
    const { login, logout } = useAuth();
    const { data: profile, isLoading, isError } = useGetProfile();

    useEffect(() => {
        if (profile) {
            login({
                id: profile.id,
                name: `${profile.first_name} ${profile.last_name}`.trim(),
                email: profile.email,
                role: profile.role,
            });
        } else if (isError) {
            logout();
        }
    }, [profile, isError, login, logout]);

    // We don't block rendering here to allow public pages or loading states in components
    return <>{children}</>;
}
