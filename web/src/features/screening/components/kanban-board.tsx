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
import { toast } from "sonner";

import { useKanbanDnd } from "../hooks/use-kanban-dnd";
import { useKanbanSelection } from "../hooks/use-kanban-selection";
import { useKanbanOutreach } from "../hooks/use-kanban-outreach";

import { useGetStages } from "../../jobs/api/use-get-stages";
import { useCreateStage } from "../../jobs/api/use-create-stage";
import { useUpdateStage } from "../../jobs/api/use-update-stage";
import { useDeleteStage } from "../../jobs/api/use-delete-stage";
import { useGetCandidates } from "../../candidates/api/use-get-candidates";
import { CandidateCard } from "../types";
import { StageModal } from "./stage-modal";
import { Stage } from "../../jobs/types/stage";

export function KanbanBoard({ jobId }: { jobId: string }) {
    const t = useTranslations("Screening");
    const [isEditMode, setIsEditMode] = useState(false);
    const [isStageModalOpen, setIsStageModalOpen] = useState(false);
    const [editingStage, setEditingStage] = useState<Stage | null>(null);

    const { data: stages = [], isLoading: isLoadingStages } = useGetStages(jobId);
    const { data: groupedCandidates = {}, isLoading: isLoadingCandidates } = useGetCandidates(jobId);

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
    } = useKanbanSelection(jobId);

    // Ensure all stages from the backend are represented in the columns object
    // even if they have no candidates yet. This is critical for drag-and-drop.
    const boardColumns = useMemo(() => {
        const cols: Record<string, CandidateCard[]> = {};
        stages.forEach((stage) => {
            cols[stage.id] = groupedCandidates[stage.id] || [];
        });
        return cols;
    }, [stages, groupedCandidates]);
    
    // Previous simplification was too aggressive - it broke dropping onto empty columns

    const {
        activeCard,
        sensors,
        collisionDetectionStrategy,
        handleDragStart,
        handleDragOver,
        handleDragEnd,
    } = useKanbanDnd({
        jobId,
        isSelectionMode,
        initialColumns: boardColumns,
        onColumnChange: handleColumnChange as (candidateId: string, card: any, sourceColumn: string, targetColumn: string) => void,
    });

    const isBoardEmpty = useMemo(() => {
        return Object.values(boardColumns).every(col => col.length === 0);
    }, [boardColumns]);

    const { mutate: createStage } = useCreateStage(jobId);
    const { mutate: updateStage } = useUpdateStage(jobId);
    const { mutate: deleteStage } = useDeleteStage(jobId);

    const handleAddStage = () => {
        setEditingStage(null);
        setIsStageModalOpen(true);
    };

    const handleEditStage = (stage: Stage) => {
        setEditingStage(stage);
        setIsStageModalOpen(true);
    };

    const onModalSubmit = (values: { title: string; is_terminal: boolean }) => {
        if (editingStage) {
            updateStage({
                stageId: editingStage.id,
                data: {
                    title: values.title,
                    is_terminal: values.is_terminal,
                    code: values.title.toLowerCase().replace(/\s+/g, "_"),
                },
            });
        } else {
            createStage({
                title: values.title,
                code: values.title.toLowerCase().replace(/\s+/g, "_"),
                position: stages.length,
                is_terminal: values.is_terminal,
            });
        }
        setIsStageModalOpen(false);
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
                            candidates={boardColumns[stage.id] || []}
                            selectedCandidateIds={selectedCandidateIds}
                            onToggleSelection={handleToggleSelection}
                            isSelectionMode={isSelectionMode}
                            isEditMode={isEditMode}
                            onEdit={() => handleEditStage(stage)}
                            onDelete={() => {
                                if (window.confirm(t("deleteStageConfirm"))) {
                                    deleteStage(stage.id);
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

            <StageModal
                isOpen={isStageModalOpen}
                onClose={() => setIsStageModalOpen(false)}
                onSubmit={onModalSubmit}
                initialValues={editingStage ? {
                    title: editingStage.title,
                    is_terminal: editingStage.is_terminal,
                } : undefined}
                title={editingStage ? t("renameStagePrompt") : t("addStage")}
            />
        </div>
    );
}
