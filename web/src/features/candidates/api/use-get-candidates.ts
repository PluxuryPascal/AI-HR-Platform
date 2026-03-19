import { useQuery } from "@tanstack/react-query";
import { useJobsStore } from "@/store/use-jobs-store";
import { CandidateCard } from "@/features/screening/types";

export function useGetCandidates(jobId?: string) {
    const boards = useJobsStore((state) => state.boards);

    return useQuery<Record<string, CandidateCard[]>>({
        queryKey: ["candidates", jobId],
        queryFn: async (): Promise<Record<string, CandidateCard[]>> => {
            if (!jobId) {
                return {};
            }
            // Artificial delay for consistency
            await new Promise((resolve) => setTimeout(resolve, 300));
            return boards[jobId] || {};
        },
        enabled: !!jobId,
        staleTime: 1000 * 60 * 5,
    });
}
