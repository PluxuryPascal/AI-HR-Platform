import { useQuery } from "@tanstack/react-query";
import { initialColumns } from "@/features/screening/utils/mock-board";
import { CandidateCard, ColumnId } from "@/features/screening/types";

// Simulated API fetch
const fetchCandidates = async (): Promise<Record<ColumnId, CandidateCard[]>> => {
    // Simulate network delay
    await new Promise((resolve) => setTimeout(resolve, 500));
    return initialColumns;
};

export function useGetCandidates() {
    return useQuery({
        queryKey: ["candidates"],
        queryFn: fetchCandidates,
        // Keep data fresh for 5 minutes since it's mock data anyway
        staleTime: 1000 * 60 * 5,
    });
}
