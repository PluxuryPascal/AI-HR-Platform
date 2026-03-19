import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Stage } from "../types/stage";
import { ApiResponse } from "@/types";

export function useGetStages(jobId?: string) {
    return useQuery({
        queryKey: ["stages", jobId],
        queryFn: async () => {
            if (!jobId) return [];
            const response = await apiClient.get<ApiResponse<Stage[]>>(`/jobs/${jobId}/stages`);
            return response.data;
        },
        enabled: !!jobId,
    });
}
