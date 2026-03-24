import { useMutation } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types";

interface CompareCandidatesRequest {
    candidate_ids: string[];
}

export type CompareResultEntry = {
    experience: string;
    skills: string[];
    salary_range: string;
    risks: string;
    summary: string;
};

export type CompareCandidatesResponse = Record<string, CompareResultEntry>;

export function useCompareCandidates(jobId: string) {
    return useMutation<CompareCandidatesResponse, Error, CompareCandidatesRequest>({
        mutationFn: async (data) => {
            const response = await apiClient.post<ApiResponse<CompareCandidatesResponse>>(
                `/jobs/${jobId}/candidates/compare`, 
                data
            );
            return response.data;
        },
    });
}
