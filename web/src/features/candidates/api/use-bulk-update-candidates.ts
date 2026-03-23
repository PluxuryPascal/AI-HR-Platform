import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { CandidateCard, ColumnId } from "@/features/screening/types";

interface BulkUpdatePayload {
    candidateIds: string[];
    newStatus: string; // This is now a Stage ID
}

export function useBulkUpdateCandidates(jobId?: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ candidateIds, newStatus }: BulkUpdatePayload) => {
            if (!jobId) throw new Error("jobId is required");
            const response = await apiClient.post(`/jobs/${jobId}/candidates/bulk-move`, {
                candidate_ids: candidateIds,
                to_stage_id: newStatus,
            });
            return response;
        },
        onMutate: async ({ candidateIds, newStatus }) => {
            await queryClient.cancelQueries({ queryKey: ["candidates", jobId] });

            const previousCandidates = queryClient.getQueryData<Record<string, CandidateCard[]>>(["candidates", jobId]);

            if (previousCandidates) {
                queryClient.setQueryData<Record<string, CandidateCard[]>>(["candidates", jobId], (old) => {
                    if (!old) return old;

                    const newData = { ...old };
                    const movedCandidates: CandidateCard[] = [];

                    // 1. Remove from all existing stages
                    Object.keys(newData).forEach((stageId) => {
                        newData[stageId] = newData[stageId].filter((candidate) => {
                            if (candidateIds.includes(candidate.id)) {
                                movedCandidates.push({ ...candidate, status: newStatus });
                                return false;
                            }
                            return true;
                        });
                    });

                    // 2. Add to target stage
                    if (!newData[newStatus]) newData[newStatus] = [];
                    newData[newStatus] = [...movedCandidates, ...newData[newStatus]];

                    return newData;
                });
            }

            return { previousCandidates };
        },
        onError: (err, newTodo, context) => {
            if (context?.previousCandidates) {
                queryClient.setQueryData(["candidates", jobId], context.previousCandidates);
            }
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["candidates", jobId] });
        },
    });
}
