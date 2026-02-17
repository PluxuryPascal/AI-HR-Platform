import { useState, useCallback } from "react";
import { useBulkUpdateCandidates } from "../api/use-bulk-update-candidates";
import { ColumnId } from "@/features/screening/types";

interface UseBulkActionsOptions {
    selectedIds: Set<string>;
    clearSelection: () => void;
}

export function useBulkActions({ selectedIds, clearSelection }: UseBulkActionsOptions) {
    const [bulkActionType, setBulkActionType] = useState<"rejection" | "invitation" | null>(null);
    const { mutate: bulkUpdate } = useBulkUpdateCandidates();

    const handleBulkReject = useCallback(() => {
        setBulkActionType("rejection");
    }, []);

    const handleBulkMove = useCallback((status: string) => {
        if (status === "interview") {
            setBulkActionType("invitation");
            return;
        }

        bulkUpdate({
            candidateIds: Array.from(selectedIds),
            newStatus: status as ColumnId
        });
        clearSelection();
    }, [selectedIds, clearSelection, bulkUpdate]);

    const handleSendEmail = useCallback((content: string) => {
        // In a real app, we would send the email here using the 'content'
        // For now, we just update the status
        const newStatus = bulkActionType === "rejection" ? "rejected" : "interview";

        bulkUpdate({
            candidateIds: Array.from(selectedIds),
            newStatus: newStatus
        });
        clearSelection();
        setBulkActionType(null);
    }, [bulkActionType, selectedIds, clearSelection, bulkUpdate]);

    const closeOutreach = useCallback(() => {
        setBulkActionType(null);
    }, []);

    return {
        bulkActionType,
        isOutreachOpen: !!bulkActionType,
        handleBulkReject,
        handleBulkMove,
        handleSendEmail,
        closeOutreach,
    };
}
