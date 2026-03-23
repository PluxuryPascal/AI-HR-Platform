import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { CandidateCard, ColumnId } from "@/features/screening/types";
import { toast } from "sonner";
import { ApiResponse } from "@/types";
import { MoveCandidateRequest } from "../types/candidate";

interface MoveCandidatePayload {
    candidateId: string;
    sourceColumnId: ColumnId;
    targetColumnId: ColumnId;
    newIndex: number;
    // Optional snapshot of the state BEFORE the optimistic update started
    // This is useful if the UI has already applied the changes (e.g., during drag-over)
    optimisticSnapshot?: Record<ColumnId, CandidateCard[]>;
}

// Real API call
const moveCandidateApi = async ({ candidateId, targetColumnId, newIndex }: MoveCandidatePayload) => {
    const response = await apiClient.post<ApiResponse<{ success: boolean }>>(
        `/candidates/${candidateId}/move`,
        {
            to_stage_id: targetColumnId,
            new_position: newIndex,
        } as MoveCandidateRequest
    );
    return response.data;
};

export function useMoveCandidate(jobId?: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: moveCandidateApi,
        onMutate: async ({ candidateId, sourceColumnId, targetColumnId, newIndex, optimisticSnapshot }) => {
            // Cancel any outgoing refetches
            await queryClient.cancelQueries({ queryKey: ["candidates", jobId] });

            // Snapshot the previous value
            const previousCandidates = optimisticSnapshot || queryClient.getQueryData<Record<ColumnId, CandidateCard[]>>(["candidates", jobId, "grouped"]);

            // If we provided a snapshot, the cache might already be updated by the drag-over handler.
            // If NOT, we need to apply the update now.
            if (!optimisticSnapshot && previousCandidates) {
                queryClient.setQueryData<Record<ColumnId, CandidateCard[]>>(["candidates", jobId, "grouped"], (old) => {
                    if (!old) return old;

                    const newColumns = { ...old };
                    const sourceList = [...newColumns[sourceColumnId]];
                    const targetList = sourceColumnId === targetColumnId ? sourceList : [...newColumns[targetColumnId]];

                    const candidateIndex = sourceList.findIndex((c) => c.id === candidateId);
                    if (candidateIndex === -1) return old;

                    const [candidate] = sourceList.splice(candidateIndex, 1);
                    targetList.splice(newIndex, 0, candidate);

                    newColumns[sourceColumnId] = sourceList;
                    newColumns[targetColumnId] = targetList;

                    return newColumns;
                });
            }

            // If optimisticSnapshot WAS provided, we assume the cache is already in the desired state due to setQueryData in dragOver.
            // We could verify/force set it here, but typically it is consistent.

            return { previousCandidates };
        },
        onError: (err, newTodo, context) => {
            // Rollback to the previous value
            if (context?.previousCandidates) {
                queryClient.setQueryData(["candidates", jobId, "grouped"], context.previousCandidates);
            }
            toast.error("Failed to move candidate. Reverting changes.");
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["candidates", jobId] });
        },
    });
}
