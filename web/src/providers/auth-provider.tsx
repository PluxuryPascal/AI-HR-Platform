"use client";

import { useEffect, type ReactNode } from "react";
import { usePathname } from "next/navigation";
import { useGetProfile } from "@/features/profile/api/use-profile";
import { useAuth } from "@/store/use-auth";

export function AuthProvider({ children }: { children: ReactNode }) {
    const pathname = usePathname();
    const { login, logout } = useAuth();
    
    // List of public paths that don't require a profile
    const isPublicPage = pathname.includes("/auth") || pathname.includes("/accept-invite");
    
    const { data: profile, isError } = useGetProfile({ enabled: !isPublicPage });

    useEffect(() => {
        if (profile) {
            login({
                id: profile.id,
                name: `${profile.first_name} ${profile.last_name}`.trim(),
                email: profile.email,
                role: profile.role,
            });
        } else if (isError && !isPublicPage) {
            logout();
        }
    }, [profile, isError, login, logout, isPublicPage]);

    // We don't block rendering here to allow public pages or loading states in components
    return <>{children}</>;
}
