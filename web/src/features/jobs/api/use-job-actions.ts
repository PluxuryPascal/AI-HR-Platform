import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Job } from "../types/job";
import { toast } from "sonner";

export function usePublishJob() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (id: string) => {
            return await apiClient.post<Job>(`/jobs/${id}/publish`);
        },
        onSuccess: (_, id) => {
            queryClient.invalidateQueries({ queryKey: ["jobs"] });
            queryClient.invalidateQueries({ queryKey: ["job", id] });
            toast.success("Job published!");
        },
    });
}

export function useCloseJob() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (id: string) => {
            return await apiClient.post<Job>(`/jobs/${id}/close`);
        },
        onSuccess: (_, id) => {
            queryClient.invalidateQueries({ queryKey: ["jobs"] });
            queryClient.invalidateQueries({ queryKey: ["job", id] });
            toast.success("Job closed!");
        },
    });
}

export function useArchiveJob() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (id: string) => {
            return await apiClient.post<Job>(`/jobs/${id}/archive`);
        },
        onSuccess: (_, id) => {
            queryClient.invalidateQueries({ queryKey: ["jobs"] });
            queryClient.invalidateQueries({ queryKey: ["job", id] });
            toast.success("Job archived!");
        },
    });
}
