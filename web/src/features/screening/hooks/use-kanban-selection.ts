import { useState, useMemo } from "react";
import { CandidateCard as CandidateCardType, ColumnId } from "../types";
import { useGetCandidates } from "../../candidates/api/use-get-candidates";
import { initialColumns } from "../utils/mock-board";

export function useKanbanSelection() {
    const { data: columns = initialColumns } = useGetCandidates();
    const [selectedCandidateIds, setSelectedCandidateIds] = useState<string[]>([]);
    const [isComparisonModalOpen, setIsComparisonModalOpen] = useState(false);
    const [isSelectionMode, setIsSelectionMode] = useState(false);

    const handleToggleSelection = (id: string) => {
        setSelectedCandidateIds((prev) =>
            prev.includes(id)
                ? prev.filter((candidateId) => candidateId !== id)
                : [...prev, id]
        );
    };

    const handleClearSelection = () => {
        setSelectedCandidateIds([]);
    };

    const toggleSelectionMode = () => {
        setIsSelectionMode(!isSelectionMode);
        if (isSelectionMode) {
            handleClearSelection();
        }
    };

    // Helper to find selected candidate objects
    const getSelectedCandidates = useMemo(() => {
        const allCandidates: CandidateCardType[] = [];
        Object.values(columns).forEach(col => {
            allCandidates.push(...col);
        });
        return allCandidates.filter(c => selectedCandidateIds.includes(c.id));
    }, [columns, selectedCandidateIds]);

    return {
        selectedCandidateIds,
        isComparisonModalOpen,
        setIsComparisonModalOpen,
        isSelectionMode,
        handleToggleSelection,
        handleClearSelection,
        toggleSelectionMode,
        getSelectedCandidates,
    };
}
