import { useMemo } from "react";
import { useGetCandidates, useGetCandidatesPaginated } from "../api/use-get-candidates";
import { CandidateCard } from "@/features/screening/types";
import { CandidateFilter } from "../types/candidate";

import { mapBackendCandidateToCard } from "../utils/mappers";

export function useCandidates(jobId: string) {
    const { data: columns, isLoading } = useGetCandidates(jobId);

    const candidates = useMemo(() => {
        if (!columns) return [];
        return Object.values(columns).flat();
    }, [columns]);

    return { data: candidates, isLoading };
}

interface UseCandidatesPaginationParams {
    jobId: string;
    page?: number;
    per_page?: number;
    filter?: CandidateFilter;
}

export function useCandidatesPagination({ jobId, page, per_page, filter }: UseCandidatesPaginationParams) {
    const { data: response, isLoading } = useGetCandidatesPaginated({ jobId, page, per_page, filter });

    const candidates = useMemo(() => {
        if (!response?.data) return [];
        
        return response.data.map(mapBackendCandidateToCard);
    }, [response]);

    return { 
        data: candidates, 
        meta: response?.meta, 
        isLoading 
    };
}
