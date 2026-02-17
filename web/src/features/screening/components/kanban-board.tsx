"use client";

import {
    DndContext,
    MeasuringStrategy,
} from "@dnd-kit/core";
import { useMemo } from "react";
import { useTranslations } from "next-intl";
import { Layers } from "lucide-react";
import { EmptyState } from "@/components/shared/empty-state";
import { BoardColumn } from "./board-column";
import { ColumnId } from "../types";

import { FloatingComparisonBar } from "./floating-comparison-bar";
import { ComparisonModal } from "./comparison-modal";
import { OutreachDrawer } from "./outreach-drawer";
import { KanbanDragOverlay } from "./kanban-drag-overlay";
import { KanbanToolbar } from "./kanban-toolbar";

import { useKanbanDnd } from "../hooks/use-kanban-dnd";
import { useKanbanSelection } from "../hooks/use-kanban-selection";
import { useKanbanOutreach } from "../hooks/use-kanban-outreach";

export function KanbanBoard() {
    const t = useTranslations("Screening");

    const { outreachState, handleColumnChange, closeOutreach } = useKanbanOutreach();

    const {
        selectedCandidateIds,
        isComparisonModalOpen,
        setIsComparisonModalOpen,
        isSelectionMode,
        handleToggleSelection,
        handleClearSelection,
        toggleSelectionMode,
        getSelectedCandidates,
    } = useKanbanSelection();

    const {
        columns,
        activeCard,
        sensors,
        collisionDetectionStrategy,
        handleDragStart,
        handleDragOver,
        handleDragEnd,
    } = useKanbanDnd({
        isSelectionMode,
        onColumnChange: handleColumnChange,
    });

    const isBoardEmpty = useMemo(() => {
        return Object.values(columns).every(col => col.length === 0);
    }, [columns]);

    if (isBoardEmpty) {
        return (
            <EmptyState
                icon={Layers}
                title={t("empty.title")}
                description={t("empty.desc")}
                className="h-[600px] border rounded-md border-dashed"
            />
        );
    }

    return (
        <div className="flex flex-col h-full gap-4">
            <KanbanToolbar
                isSelectionMode={isSelectionMode}
                onToggleSelectionMode={toggleSelectionMode}
            />

            <DndContext
                sensors={sensors}
                collisionDetection={collisionDetectionStrategy}
                onDragStart={handleDragStart}
                onDragOver={handleDragOver}
                onDragEnd={handleDragEnd}
                measuring={{
                    droppable: {
                        strategy: MeasuringStrategy.Always,
                    },
                }}
            >
                <div className="flex h-full gap-4 overflow-x-auto pb-4 relative">
                    {(Object.keys(columns) as ColumnId[]).map((columnId) => (
                        <BoardColumn
                            key={columnId}
                            id={columnId}
                            candidates={columns[columnId]}
                            selectedCandidateIds={selectedCandidateIds}
                            onToggleSelection={handleToggleSelection}
                            isSelectionMode={isSelectionMode}
                        />
                    ))}
                </div>
                <KanbanDragOverlay activeCard={activeCard} />

                {isSelectionMode && (
                    <FloatingComparisonBar
                        selectedCount={selectedCandidateIds.length}
                        onCompare={() => setIsComparisonModalOpen(true)}
                        onClear={handleClearSelection}
                    />
                )}

                <ComparisonModal
                    isOpen={isComparisonModalOpen}
                    onClose={() => setIsComparisonModalOpen(false)}
                    candidates={getSelectedCandidates}
                />

                <OutreachDrawer
                    isOpen={outreachState.isOpen}
                    onClose={closeOutreach}
                    candidate={outreachState.candidate}
                    type={outreachState.type}
                />
            </DndContext>
        </div>
    );
}
