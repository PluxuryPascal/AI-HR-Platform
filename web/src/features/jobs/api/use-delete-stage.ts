import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";

export function useDeleteStage(jobId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (stageId: string) => {
            await apiClient.delete(`/jobs/${jobId}/stages/${stageId}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["stages", jobId] });
            toast.success("Stage deleted successfully");
        },
        onError: (error: Error) => {
            toast.error(`Failed to delete stage: ${error.message}`);
        },
    });
}
