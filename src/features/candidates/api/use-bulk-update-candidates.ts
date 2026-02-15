import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CandidateCard, ColumnId } from "@/features/screening/types";

interface BulkUpdatePayload {
    candidateIds: string[];
    newStatus: ColumnId;
}

export function useBulkUpdateCandidates() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ candidateIds, newStatus }: BulkUpdatePayload) => {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 800));
            return { candidateIds, newStatus };
        },
        onMutate: async ({ candidateIds, newStatus }) => {
            await queryClient.cancelQueries({ queryKey: ["candidates"] });

            const previousCandidates = queryClient.getQueryData<Record<ColumnId, CandidateCard[]>>(["candidates"]);

            queryClient.setQueryData<Record<ColumnId, CandidateCard[]>>(["candidates"], (old) => {
                if (!old) return old;

                const newData = { ...old };
                const movedCandidates: CandidateCard[] = [];

                // Find and remove candidates from their current columns
                (Object.keys(newData) as ColumnId[]).forEach((columnId) => {
                    newData[columnId] = newData[columnId].filter((candidate) => {
                        if (candidateIds.includes(candidate.id)) {
                            movedCandidates.push({ ...candidate, status: newStatus }); // Update status
                            return false;
                        }
                        return true;
                    });
                });

                // Add candidates to the new column
                if (newData[newStatus]) {
                    newData[newStatus] = [...movedCandidates, ...newData[newStatus]];
                } else {
                    // Fallback if column needed to be created, though it should exist
                    newData[newStatus] = movedCandidates;
                }

                return newData;
            });

            return { previousCandidates };
        },
        onError: (err, newTodo, context) => {
            if (context?.previousCandidates) {
                queryClient.setQueryData(["candidates"], context.previousCandidates);
            }
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["candidates"] });
        },
    });
}
