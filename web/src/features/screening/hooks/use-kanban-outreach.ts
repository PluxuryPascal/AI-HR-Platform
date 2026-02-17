import { useState } from "react";
import { CandidateCard as CandidateCardType, ColumnId } from "../types";

interface OutreachState {
    isOpen: boolean;
    candidate: CandidateCardType | null;
    type: "rejection" | "invitation";
}

export function useKanbanOutreach() {
    const [outreachState, setOutreachState] = useState<OutreachState>({
        isOpen: false,
        candidate: null,
        type: "rejection",
    });

    const handleColumnChange = (
        _candidateId: string,
        card: CandidateCardType,
        _sourceColumn: ColumnId,
        targetColumn: ColumnId,
    ) => {
        if (targetColumn === "rejected") {
            setOutreachState({
                isOpen: true,
                candidate: card,
                type: "rejection",
            });
        } else if (targetColumn === "interview") {
            setOutreachState({
                isOpen: true,
                candidate: card,
                type: "invitation",
            });
        }
    };

    const closeOutreach = () => {
        setOutreachState(prev => ({ ...prev, isOpen: false }));
    };

    return {
        outreachState,
        handleColumnChange,
        closeOutreach,
    };
}
