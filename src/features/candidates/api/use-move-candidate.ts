import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CandidateCard, ColumnId } from "@/features/screening/types";
import { toast } from "sonner";

interface MoveCandidatePayload {
    candidateId: string;
    sourceColumnId: ColumnId;
    targetColumnId: ColumnId;
    newIndex: number;
    // Optional snapshot of the state BEFORE the optimistic update started
    // This is useful if the UI has already applied the changes (e.g., during drag-over)
    optimisticSnapshot?: Record<ColumnId, CandidateCard[]>;
}

// Simulated API call
const moveCandidateApi = async ({ optimisticSnapshot, ...payload }: MoveCandidatePayload) => {
    // In a real app, payload would be sanitized
    await new Promise((resolve) => setTimeout(resolve, 500));
    return { success: true };
};

export function useMoveCandidate() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: moveCandidateApi,
        onMutate: async ({ candidateId, sourceColumnId, targetColumnId, newIndex, optimisticSnapshot }) => {
            // Cancel any outgoing refetches
            await queryClient.cancelQueries({ queryKey: ["candidates"] });

            // Snapshot the previous value
            // If optimisticSnapshot is provided, use it (assumes UI is already updated)
            // Otherwise, fetch current cache (assumes UI is waiting for this update)
            const previousCandidates = optimisticSnapshot || queryClient.getQueryData<Record<ColumnId, CandidateCard[]>>(["candidates"]);

            // If we provided a snapshot, the cache might already be updated by the drag-over handler.
            // If NOT, we need to apply the update now.
            if (!optimisticSnapshot && previousCandidates) {
                queryClient.setQueryData<Record<ColumnId, CandidateCard[]>>(["candidates"], (old) => {
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
                queryClient.setQueryData(["candidates"], context.previousCandidates);
            }
            toast.error("Failed to move candidate. Reverting changes.");
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["candidates"] });
        },
    });
}
