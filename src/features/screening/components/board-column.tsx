"use client";

import { useDroppable } from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { ColumnId, CandidateCard as CandidateCardType } from "../types";
import { CandidateCard } from "./candidate-card";
import { useTranslations } from "next-intl";

interface BoardColumnProps {
    id: ColumnId;
    candidates: CandidateCardType[];
    selectedCandidateIds: string[];
    onToggleSelection: (id: string) => void;
    isSelectionMode: boolean;
}

export function BoardColumn({ id, candidates, selectedCandidateIds, onToggleSelection, isSelectionMode }: BoardColumnProps) {
    const t = useTranslations('Screening');
    const { setNodeRef } = useDroppable({
        id: id,
    });

    return (
        <div
            ref={setNodeRef}
            className="flex flex-col bg-muted/50 rounded-lg p-4 min-w-[300px] h-full"
        >
            <div className="flex items-center justify-between mb-4">
                <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">
                    {t(`columns.${id}`)}
                </h3>
                <span className="text-xs bg-background px-2 py-1 rounded-full shadow-sm">
                    {candidates.length}
                </span>
            </div>

            <div className="flex-1 flex flex-col gap-3 min-h-[100px]">
                <SortableContext
                    items={candidates.map((c) => c.id)}
                    strategy={verticalListSortingStrategy}
                >
                    {candidates.map((candidate) => (
                        <CandidateCard
                            key={candidate.id}
                            candidate={candidate}
                            isSelected={selectedCandidateIds.includes(candidate.id)}
                            onToggleSelection={onToggleSelection}
                            isSelectionMode={isSelectionMode}
                        />
                    ))}
                </SortableContext>
                {candidates.length === 0 && (
                    <div className="flex items-center justify-center h-20 border-2 border-dashed rounded-md border-muted-foreground/20 text-muted-foreground text-xs">
                        {t('empty')}
                    </div>
                )}
            </div>
        </div>
    );
}
