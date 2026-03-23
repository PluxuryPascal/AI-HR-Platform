import { create } from 'zustand';

interface User {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    avatar?: string;
    role?: string;
}

interface AuthState {
    isAuthenticated: boolean;
    user: User | null;
    login: (userData: User) => void;
    logout: () => void;
}

export const useAuth = create<AuthState>((set) => ({
    isAuthenticated: false,
    user: null,
    login: (userData) => set({ isAuthenticated: true, user: userData }),
    logout: () => set({ isAuthenticated: false, user: null }),
}));
