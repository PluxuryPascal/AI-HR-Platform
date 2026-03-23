import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { CreateStageRequest, Stage } from "../types/stage";
import { ApiResponse } from "@/types";
import { toast } from "sonner";

export function useCreateStage(jobId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: CreateStageRequest) => {
            const response = await apiClient.post<ApiResponse<Stage>>(
                `/jobs/${jobId}/stages`,
                data
            );
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["stages", jobId] });
            toast.success("Stage created successfully");
        },
        onError: (error: Error) => {
            toast.error(`Failed to create stage: ${error.message}`);
        },
    });
}
