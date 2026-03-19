import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Job } from "../types/job";
import { toast } from "sonner";

export interface CreateJobPayload {
    title: string;
    department_id?: string;
    work_format: "remote" | "office" | "hybrid";
    description?: string;
    salary_min?: number;
    salary_max?: number;
    currency?: string;
}

export function useCreateJob() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (payload: CreateJobPayload) => {
            return await apiClient.post<Job>("/jobs", payload);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["jobs"] });
            toast.success("Вакансия успешно создана");
        },
        onError: (error) => {
            console.error("Create job error:", error);
            toast.error("Ошибка при создании вакансии");
        },
    });
}
