import { useMutation } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { RegisterRequest } from "../types";

export function useRegister() {
    return useMutation({
        mutationFn: async (data: RegisterRequest) => {
            const response = await apiClient.post("/auth/register", data);
            return response;
        },
    });
}

// Add useLogin as well for completeness if needed later
export interface LoginRequest {
    email: string;
    password: string;
}

export function useLogin() {
    return useMutation({
        mutationFn: async (data: LoginRequest) => {
            const response = await apiClient.post("/auth/login", data);
            return response;
        },
    });
}
