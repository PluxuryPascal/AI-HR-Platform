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
import { useGetCandidates } from "../../candidates/api/use-get-candidates";
import { CandidateCard } from "../types";

export function KanbanBoard({ jobId }: { jobId: string }) {
    const t = useTranslations("Screening");
    const [isEditMode, setIsEditMode] = useState(false);

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
    } = useKanbanSelection();

    // Map backend candidates to CandidateCard UI format
    const boardColumns = useMemo(() => {
        const mapped: Record<string, CandidateCard[]> = {};
        stages.forEach(stage => {
            const candidates = groupedCandidates[stage.id] || [];
            mapped[stage.id] = (candidates as any[]).map(c => ({
                id: c.id,
                name: `${c.first_name || ""} ${c.last_name || ""}`.trim() || t("unknownCandidate") || "Unknown",
                role: t("candidateRole") || "Candidate",
                score: 0,
                avatarUrl: `https://api.dicebear.com/7.x/avataaars/svg?seed=${c.id}`,
                email: c.email || "-",
                appliedDate: new Date(c.created_at).toLocaleDateString(),
                status: c.parsing_status,
                matchSummary: "",
                scoreBreakdown: []
            }));
        });
        return mapped;
    }, [stages, groupedCandidates, t]);

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

    const handleAddStage = () => {
        // In a real app, this would call useCreateStage mutation
        toast.info("Эта фича будет доступна скоро");
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
                            onRename={(newTitle) => toast.info("Rename feature coming soon")}
                            onDelete={() => {
                                if (window.confirm(t("deleteStageConfirm"))) {
                                    toast.info("Delete feature coming soon");
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
