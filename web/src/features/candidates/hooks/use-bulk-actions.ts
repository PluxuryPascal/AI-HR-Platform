import { useState, useCallback } from "react";
import { useBulkUpdateCandidates } from "../api/use-bulk-update-candidates";
import { useBulkDeleteCandidates } from "../api/use-bulk-delete-candidates";
import { ColumnId } from "@/features/screening/types";

interface UseBulkActionsOptions {
    selectedIds: Set<string>;
    clearSelection: () => void;
    jobId?: string;
    stages?: { id: string; code: string }[];
}

export function useBulkActions({ selectedIds, clearSelection, jobId, stages }: UseBulkActionsOptions) {
    const [bulkActionType, setBulkActionType] = useState<"rejection" | "invitation" | null>(null);
    const { mutate: bulkUpdate } = useBulkUpdateCandidates(jobId);
    const { mutate: bulkDelete } = useBulkDeleteCandidates(jobId);

    const handleBulkReject = useCallback(() => {
        setBulkActionType("rejection");
    }, []);

    const handleBulkMove = useCallback((stageId: string) => {
        bulkUpdate({
            candidateIds: Array.from(selectedIds),
            newStatus: stageId
        });
        clearSelection();
    }, [selectedIds, clearSelection, bulkUpdate]);

    const handleBulkDelete = useCallback(() => {
        if (window.confirm("Are you sure you want to delete these candidates?")) {
            bulkDelete({
                candidateIds: Array.from(selectedIds)
            });
            clearSelection();
        }
    }, [selectedIds, clearSelection, bulkDelete]);

    const handleSendEmail = useCallback((content: string) => {
        // Find the correct stage ID based on the action type
        const targetCode = bulkActionType === "rejection" ? "reject" : "interview";
        const targetStage = stages?.find(s => s.code === targetCode);

        if (!targetStage) {
            console.error(`Stage with code ${targetCode} not found`);
            return;
        }

        bulkUpdate({
            candidateIds: Array.from(selectedIds),
            newStatus: targetStage.id
        });
        clearSelection();
        setBulkActionType(null);
    }, [bulkActionType, selectedIds, clearSelection, bulkUpdate, stages]);

    const closeOutreach = useCallback(() => {
        setBulkActionType(null);
    }, []);

    return {
        bulkActionType,
        isOutreachOpen: !!bulkActionType,
        handleBulkReject,
        handleBulkDelete,
        handleBulkMove,
        handleSendEmail,
        closeOutreach,
    };
}
