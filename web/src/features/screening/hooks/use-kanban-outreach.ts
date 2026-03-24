import { useState } from "react";
import confetti from "canvas-confetti";
import { CandidateCard as CandidateCardType, ColumnId } from "../types";
import { Stage } from "../../jobs/types/stage";

interface OutreachState {
    isOpen: boolean;
    candidate: CandidateCardType | null;
    type: "rejection" | "invitation";
}

export function useKanbanOutreach(stages: Stage[]) {
    const [outreachState, setOutreachState] = useState<OutreachState>({
        isOpen: false,
        candidate: null,
        type: "rejection",
    });

    const handleColumnChange = (
        _candidateId: string,
        card: CandidateCardType,
        _sourceColumn: ColumnId,
        targetColumnId: ColumnId,
    ) => {
        const targetStage = stages.find(s => s.id === targetColumnId);
        if (!targetStage) return;

        if (targetStage.is_rejection) {
            setOutreachState({
                isOpen: true,
                candidate: card,
                type: "rejection",
            });
        } else if (targetStage.is_interview) {
            confetti({
                particleCount: 150,
                spread: 70,
                origin: { y: 0.6 },
                colors: ["#2563eb", "#4f46e5", "#818cf8"],
            });
            setOutreachState({
                isOpen: true,
                candidate: card,
                type: "invitation",
            });
        }
    };

    const closeOutreach = () => {
        setOutreachState((prev: OutreachState) => ({ ...prev, isOpen: false }));
    };

    return {
        outreachState,
        handleColumnChange,
        closeOutreach,
    };
}
