import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Candidate as BackendCandidate, CandidateFilter } from "../types/candidate";
import { PaginatedResponse } from "@/types";
import { CandidateCard } from "@/features/screening/types";

import { mapBackendCandidateToCard } from "../utils/mappers";

// Original hook for Kanban board (grouped by stage)
export function useGetCandidates(jobId?: string) {
    return useQuery({
        queryKey: ["candidates", jobId, "grouped"],
        queryFn: async () => {
            if (!jobId) {
                return {} as Record<string, CandidateCard[]>;
            }

            const response = await apiClient.post<PaginatedResponse<BackendCandidate>>(
                `/jobs/${jobId}/candidates/list`,
                {
                    pagination: { page: 1, per_page: 500 }, // Large limit for board
                }
            );

            const grouped: Record<string, CandidateCard[]> = {};
            
            response.data.forEach((c) => {
                const stageId = c.stage_id || "new"; 
                if (!grouped[stageId]) grouped[stageId] = [];
                grouped[stageId].push(mapBackendCandidateToCard(c));
            });

            return grouped;
        },
        enabled: !!jobId,
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
    });
}

interface GetCandidatesParams {
    jobId: string;
    page?: number;
    per_page?: number;
    filter?: CandidateFilter;
}

// New hook for Table view (paginated, sorted, filtered)
export function useGetCandidatesPaginated({ jobId, page = 1, per_page = 10, filter }: GetCandidatesParams) {
    return useQuery({
        queryKey: ["candidates", "paginated", jobId, page, per_page, filter],
        queryFn: async () => {
            if (!jobId) {
                return null;
            }

            const response = await apiClient.post<PaginatedResponse<BackendCandidate>>(
                `/jobs/${jobId}/candidates/list`,
                {
                    pagination: { page, per_page },
                    filter,
                }
            );

            return response;
        },
        enabled: !!jobId,
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
    });
}
