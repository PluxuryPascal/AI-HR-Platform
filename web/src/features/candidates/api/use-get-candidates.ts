import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Candidate } from "../types/candidate";
import { PaginatedResponse } from "@/types";

export function useGetCandidates(jobId?: string) {
    return useQuery<Record<string, Candidate[]>>({
        queryKey: ["candidates", jobId],
        queryFn: async (): Promise<Record<string, Candidate[]>> => {
            if (!jobId) {
                return {};
            }

            const response = await apiClient.post<PaginatedResponse<Candidate>>(
                `/jobs/${jobId}/candidates/list`,
                {
                    pagination: { page: 1, per_page: 100 },
                }
            );

            // Group candidates by their stage
            // Note: Currently backend doesn't return stage_id in Candidate 
            // but we can assume they are grouped somehow or we need to fetch stages first.
            // Looking at CandidateHandler.GetCandidate, it returns stage_id.
            // Let's assume for now we might need to adjust this once we see the list response.
            
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
