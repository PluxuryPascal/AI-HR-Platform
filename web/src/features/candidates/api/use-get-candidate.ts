import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { CandidateDetailResponse } from "../types/candidate";
import { ApiResponse } from "@/types";

export function useGetCandidate(id: string) {
    return useQuery<CandidateDetailResponse>({
        queryKey: ["candidate", id],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<CandidateDetailResponse>>(`/candidates/${id}`);
            return response.data;
        },
        enabled: !!id,
    });
}
