import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Department, CreateDepartmentPayload } from "../types";
import { toast } from "sonner";
import { ApiResponse } from "@/types";

export function useGetDepartments() {
    return useQuery<Department[]>({
        queryKey: ["departments"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<Department[]>>("/departments");
            return response.data;
        },
    });
}

export function useCreateDepartment() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (payload: CreateDepartmentPayload) => {
            const response = await apiClient.post<ApiResponse<Department>>("/departments", payload);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["departments"] });
            toast.success("Отдел успешно создан");
        },
        onError: () => {
            toast.error("Ошибка при создании отдела");
        },
    });
}

export function useDeleteDepartment() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            const response = await apiClient.delete(`/departments/${id}`);
            return response;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["departments"] });
            toast.success("Отдел успешно удален");
        },
        onError: () => {
            toast.error("Ошибка при удалении отдела");
        },
    });
}
