"use client";

import {
    DndContext,
    DragOverlay,
    closestCorners,
    KeyboardSensor,
    PointerSensor,
    useSensor,
    useSensors,
    DragStartEvent,
    DragOverEvent,
    DragEndEvent,
    defaultDropAnimationSideEffects,
    DropAnimation,
    CollisionDetection,
    pointerWithin,
    rectIntersection,
    getFirstCollision,
    MeasuringStrategy,
} from "@dnd-kit/core";
import {
    arrayMove,
    sortableKeyboardCoordinates,
} from "@dnd-kit/sortable";
import { useState, useCallback, useMemo } from "react";
import { useTranslations } from "next-intl";
import { Layers } from "lucide-react";
import { EmptyState } from "@/components/shared/empty-state";
import { BoardColumn } from "./board-column";
import { CandidateCard } from "./candidate-card";
import { initialColumns } from "../utils/mock-board";
import { CandidateCard as CandidateCardType, ColumnId } from "../types";

import { FloatingComparisonBar } from "./floating-comparison-bar";
import { ComparisonModal } from "./comparison-modal";

import { Button } from "@/components/ui/button";

export function KanbanBoard() {
    const t = useTranslations("Screening");
    const [columns, setColumns] = useState(initialColumns);
    const [activeCard, setActiveCard] = useState<CandidateCardType | null>(null);
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
    }

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

    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: {
                distance: 5,
            },
        }),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );

    const findContainer = (id: string): ColumnId | undefined => {
        if (id in columns) {
            return id as ColumnId;
        }

        return (Object.keys(columns) as ColumnId[]).find((key) =>
            columns[key].find((c) => c.id === id)
        );
    };

    const handleDragStart = (event: DragStartEvent) => {
        if (isSelectionMode) return;
        const { active } = event;
        const { id } = active;
        const activeColumn = findContainer(id as string);
        if (activeColumn) {
            const card = columns[activeColumn].find((c) => c.id === id);
            if (card) {
                setActiveCard(card);
            }
        }
    };

    const handleDragOver = (event: DragOverEvent) => {
        if (isSelectionMode) return;
        const { active, over } = event;
        const { id } = active;
        const overId = over?.id;

        if (!overId || active.id === overId) return;

        const activeContainer = findContainer(id as string);
        const overContainer = findContainer(overId as string);

        if (
            !activeContainer ||
            !overContainer ||
            activeContainer === overContainer
        ) {
            return;
        }

        setColumns((prev) => {
            const activeItems = prev[activeContainer];
            const overItems = prev[overContainer];
            const activeIndex = activeItems.findIndex((item) => item.id === id);
            const overIndex = overItems.findIndex((item) => item.id === overId);

            let newIndex;

            if (overId in prev) {
                // We're over a container
                newIndex = overItems.length + 1;
            } else {
                const isBelowOverItem =
                    over &&
                    active.rect.current.translated &&
                    active.rect.current.translated.top >
                    over.rect.top + over.rect.height;

                const modifier = isBelowOverItem ? 1 : 0;

                newIndex =
                    overIndex >= 0 ? overIndex + modifier : overItems.length + 1;
            }

            return {
                ...prev,
                [activeContainer]: [
                    ...prev[activeContainer].filter((item) => item.id !== active.id),
                ],
                [overContainer]: [
                    ...prev[overContainer].slice(0, newIndex),
                    activeItems[activeIndex],
                    ...prev[overContainer].slice(
                        newIndex,
                        prev[overContainer].length
                    ),
                ],
            };
        });
    };

    const handleDragEnd = (event: DragEndEvent) => {
        if (isSelectionMode) return;
        const { active, over } = event;
        const { id } = active;
        const overId = over?.id;

        if (!overId) {
            setActiveCard(null);
            return;
        }

        const activeContainer = findContainer(id as string);
        const overContainer = findContainer(overId as string);

        if (
            activeContainer &&
            overContainer &&
            activeContainer === overContainer
        ) {
            const activeIndex = columns[activeContainer].findIndex(
                (item) => item.id === id
            );
            const overIndex = columns[overContainer].findIndex(
                (item) => item.id === overId
            );

            if (activeIndex !== overIndex) {
                setColumns((prev) => ({
                    ...prev,
                    [activeContainer]: arrayMove(
                        prev[activeContainer],
                        activeIndex,
                        overIndex
                    ),
                }));
            }
        }

        setActiveCard(null);
    };

    const collisionDetectionStrategy: CollisionDetection = useCallback(
        (args) => {
            if (activeCard && activeCard.id === args.active.id) {
                return closestCorners(args);
            }

            // First, attempt to detect collisions with pointers within the container
            const pointerIntersections = pointerWithin(args);
            const intersections =
                pointerIntersections.length > 0
                    ? pointerIntersections
                    : rectIntersection(args);

            let overId = getFirstCollision(intersections, "id");

            if (overId != null) {
                if (overId in columns) {
                    const containerIntersections = rectIntersection({
                        ...args,
                        droppableContainers: args.droppableContainers.filter(
                            (container) => container.id in columns
                        ),
                    });

                    if (containerIntersections.length > 0) {
                        overId = containerIntersections[0].id;
                    }
                }

                return [{ id: overId }];
            }

            return closestCorners(args);
        },
        [activeCard, columns]
    );




    const dropAnimation: DropAnimation = {
        sideEffects: defaultDropAnimationSideEffects({
            styles: {
                active: {
                    opacity: "0.5",
                },
            },
        }),

    };

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
            <div className="flex justify-end px-2">
                <Button
                    variant={isSelectionMode ? "secondary" : "outline"}
                    onClick={toggleSelectionMode}
                    size="sm"
                >
                    {isSelectionMode ? t("bar.cancelSelection") : t("bar.selectCandidates")}
                </Button>
            </div>

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
                <DragOverlay dropAnimation={dropAnimation}>
                    {activeCard ? <CandidateCard candidate={activeCard} isOverlay /> : null}
                </DragOverlay>

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
            </DndContext>
        </div>
    );
}
