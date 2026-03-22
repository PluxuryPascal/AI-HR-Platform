import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { UpdateJobRequest, Job } from "../types/job";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

export function useUpdateJob(id: string) {
    const queryClient = useQueryClient();
    const t = useTranslations("JobWizard");

    return useMutation({
        mutationFn: async (data: UpdateJobRequest) => {
            return await apiClient.patch<Job>(`/jobs/${id}`, data);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["jobs"] });
            queryClient.invalidateQueries({ queryKey: ["job", id] });
            toast.success(t("successUpdate") || "Job updated successfully!");
        },
        onError: (error: any) => {
            toast.error(error?.response?.data?.message || "Failed to update job");
        },
    });
}
