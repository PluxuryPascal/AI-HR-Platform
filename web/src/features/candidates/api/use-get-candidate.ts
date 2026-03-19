import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { CandidateDetailResponse } from "../types/candidate";

export function useGetCandidate(id: string) {
    return useQuery<CandidateDetailResponse>({
        queryKey: ["candidate", id],
        queryFn: async () => {
            const response = await apiClient.get<CandidateDetailResponse>(`/candidates/${id}`);
            return response;
        },
        enabled: !!id,
    });
}
