import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { UpdateStageRequest } from "../types/stage";
import { toast } from "sonner";

export function useUpdateStage(jobId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ stageId, data }: { stageId: string; data: UpdateStageRequest }) => {
            const response = await apiClient.patch(
                `/jobs/${jobId}/stages/${stageId}`,
                data
            );
            return response;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["stages", jobId] });
            toast.success("Stage updated successfully");
        },
        onError: (error: Error) => {
            toast.error(`Failed to update stage: ${error.message}`);
        },
    });
}
