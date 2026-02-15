import { create } from 'zustand';

interface User {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    role?: "owner" | "recruiter" | "manager";
}

interface AuthState {
    isAuthenticated: boolean;
    user: User | null;
    login: (userData: User) => void;
    logout: () => void;
}

export const useAuth = create<AuthState>((set) => ({
    isAuthenticated: true, // Default to true for testing
    user: {
        id: "1",
        name: "Test User",
        email: "test@example.com",
        role: "owner", // Default to owner
    },
    login: (userData) => set({ isAuthenticated: true, user: userData }),
    logout: () => set({ isAuthenticated: false, user: null }),
}));
