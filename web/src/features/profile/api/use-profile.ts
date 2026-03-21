import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types";

export interface UserProfile {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    role: string;
    team_id: string;
    team_name: string;
}

export interface UpdateProfileRequest {
    first_name: string;
    last_name: string;
    email: string;
}

export interface UpdatePasswordRequest {
    current_password: string;
    new_password: string;
}

export function useGetProfile() {
    return useQuery({
        queryKey: ["user-profile"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<UserProfile>>("/auth/profile");
            return response?.data;
        },
    });
}

export function useUpdateProfile() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: UpdateProfileRequest) => {
            const response = await apiClient.patch<ApiResponse<void>>("/auth/profile", data);
            return response?.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["user-profile"] });
        },
    });
}

export function useUpdatePassword() {
    return useMutation({
        mutationFn: async (data: UpdatePasswordRequest) => {
            const response = await apiClient.patch<ApiResponse<void>>("/auth/password", data);
            return response?.data;
        },
    });
}
