import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Candidate } from "../types/candidate";
import { PaginatedResponse } from "@/types";

export function useGetCandidates(jobId?: string) {
    return useQuery({
        queryKey: ["candidates", jobId],
        queryFn: async () => {
            if (!jobId) {
                return { data: [] };
            }

            return await apiClient.post<PaginatedResponse<Candidate>>(
                `/jobs/${jobId}/candidates/list`,
                {
                    pagination: { page: 1, per_page: 100 },
                }
            );
        },
        select: (response) => {
            const grouped: Record<string, Candidate[]> = {};
            
            response.data.forEach((c) => {
                const stageId = c.stage_id || "new"; 
                if (!grouped[stageId]) grouped[stageId] = [];
                grouped[stageId].push(c);
            });

            return grouped;
        },
        enabled: !!jobId,
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
    });
}
