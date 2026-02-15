"use client";

import { useRef } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
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

    const parentRef = useRef<HTMLDivElement>(null);

    const rowVirtualizer = useVirtualizer({
        count: candidates.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 120, // Estimate row height for accurate scrollbar
        overscan: 5,
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

            <div
                ref={parentRef}
                className="flex-1 overflow-y-auto min-h-[100px]"
            >
                <SortableContext
                    items={candidates.map((c) => c.id)}
                    strategy={verticalListSortingStrategy}
                >
                    <div
                        style={{
                            height: `${rowVirtualizer.getTotalSize()}px`,
                            width: '100%',
                            position: 'relative',
                        }}
                    >
                        {rowVirtualizer.getVirtualItems().map((virtualItem) => {
                            const candidate = candidates[virtualItem.index];
                            if (!candidate) return null;

                            return (
                                <div
                                    key={candidate.id}
                                    style={{
                                        position: 'absolute',
                                        top: 0,
                                        left: 0,
                                        width: '100%',
                                        height: `${virtualItem.size}px`,
                                        transform: `translateY(${virtualItem.start}px)`,
                                    }}
                                >
                                    <CandidateCard
                                        candidate={candidate}
                                        isSelected={selectedCandidateIds.includes(candidate.id)}
                                        onToggleSelection={onToggleSelection}
                                        isSelectionMode={isSelectionMode}
                                    />
                                </div>
                            );
                        })}
                    </div>
                </SortableContext>

                {candidates.length === 0 && (
                    <div className="flex items-center justify-center h-20 border-2 border-dashed rounded-md border-muted-foreground/20 text-muted-foreground text-xs mt-4">
                        {t('columnEmpty')}
                    </div>
                )}
            </div>
        </div>
    );
}
