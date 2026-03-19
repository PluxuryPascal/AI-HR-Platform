"use client";

import {
    DndContext,
    MeasuringStrategy,
} from "@dnd-kit/core";
import { useMemo, useState } from "react";
import { useTranslations } from "next-intl";
import { Layers } from "lucide-react";
import { BoardColumn } from "./board-column";
import { CandidateCard as CandidateCardComponent } from "./candidate-card";

import { FloatingComparisonBar } from "./floating-comparison-bar";
import { ComparisonModal } from "./comparison-modal";
import { OutreachDrawer } from "./outreach-drawer";
import { KanbanDragOverlay } from "./kanban-drag-overlay";
import { KanbanToolbar } from "./kanban-toolbar";

import { useKanbanDnd } from "../hooks/use-kanban-dnd";
import { useKanbanSelection } from "../hooks/use-kanban-selection";
import { useKanbanOutreach } from "../hooks/use-kanban-outreach";
import { useJobsStore } from "@/store/use-jobs-store";

export function KanbanBoard({ jobId }: { jobId: string }) {
    const t = useTranslations("Screening");
    const [isEditMode, setIsEditMode] = useState(false);

    const { jobs, addStage, removeStage, renameStage, reorderStages } = useJobsStore();
    const job = jobs.find(j => j.id === jobId);
    const stages = job?.stages || [];

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
        jobId,
        isSelectionMode,
        onColumnChange: handleColumnChange as (candidateId: string, card: any, sourceColumn: string, targetColumn: string) => void,
    });

    const isBoardEmpty = useMemo(() => {
        return Object.values(columns).every(col => (col as any[]).length === 0);
    }, [columns]);

    const handleAddStage = () => {
        const title = window.prompt(t("addStagePrompt"));
        if (title) {
            addStage(jobId, title);
        }
    };

    return (
        <div className="flex flex-col h-full gap-4">
            <KanbanToolbar
                isSelectionMode={isSelectionMode}
                onToggleSelectionMode={toggleSelectionMode}
                isEditMode={isEditMode}
                onToggleEditMode={() => setIsEditMode(!isEditMode)}
                onAddStage={handleAddStage}
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
                    {stages.map((stage) => (
                        <BoardColumn
                            key={stage.id}
                            id={stage.id}
                            title={stage.title}
                            candidates={columns[stage.id] || []}
                            selectedCandidateIds={selectedCandidateIds}
                            onToggleSelection={handleToggleSelection}
                            isSelectionMode={isSelectionMode}
                            isEditMode={isEditMode}
                            onRename={(newTitle) => renameStage(jobId, stage.id, newTitle)}
                            onDelete={() => {
                                if (window.confirm(t("deleteStageConfirm"))) {
                                    removeStage(jobId, stage.id);
                                }
                            }}
                        />
                    ))}
                    {isEditMode && (
                        <button
                            onClick={handleAddStage}
                            className="flex flex-col items-center justify-center min-w-[320px] h-full border-2 border-dashed rounded-lg hover:bg-accent/50 transition-colors gap-2 text-muted-foreground mr-4"
                        >
                            <Layers className="w-8 h-8" />
                            <span className="font-medium">{t("addStage")}</span>
                        </button>
                    )}
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
