import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";

export function useReorderStages(jobId: string) {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (stageIds: string[]) => {
            return await apiClient.put(`/jobs/${jobId}/pipeline/order`, { stage_ids: stageIds });
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["stages", jobId] });
            toast.success("Stages reordered!");
        },
        onError: (err) => {
            toast.error("Failed to reorder stages");
            console.error(err);
        }
    });
}
