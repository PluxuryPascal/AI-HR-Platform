import { useState, useCallback, useRef } from "react";
import {
    closestCorners,
    KeyboardSensor,
    PointerSensor,
    useSensor,
    useSensors,
    DragStartEvent,
    DragOverEvent,
    DragEndEvent,
    CollisionDetection,
    pointerWithin,
    rectIntersection,
    getFirstCollision,
} from "@dnd-kit/core";
import {
    arrayMove,
    sortableKeyboardCoordinates,
} from "@dnd-kit/sortable";
import { useQueryClient } from "@tanstack/react-query";
import { useMoveCandidate } from "../../candidates/api/use-move-candidate";
import { CandidateCard as CandidateCardType } from "../types";
import { fireOfferConfetti } from "@/lib/confetti";

interface UseKanbanDndOptions {
    jobId: string;
    isSelectionMode: boolean;
    initialColumns: Record<string, CandidateCardType[]>;
    onColumnChange?: (candidateId: string, card: CandidateCardType, sourceColumn: string, targetColumn: string) => void;
}

export function useKanbanDnd({ jobId, isSelectionMode, initialColumns: columns, onColumnChange }: UseKanbanDndOptions) {
    const queryClient = useQueryClient();
    const { mutate: moveCandidateApi } = useMoveCandidate(jobId);

    const previousColumnsRef = useRef<Record<string, CandidateCardType[]> | null>(null);

    const [activeCard, setActiveCard] = useState<CandidateCardType | null>(null);
    const [activeColumn, setActiveColumn] = useState<string | null>(null);

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

    const findContainer = (id: string): string | undefined => {
        if (id in columns) {
            return id;
        }

        return Object.keys(columns).find((key) => {
            const col = columns[key];
            return Array.isArray(col) && col.find((c) => c.id === id);
        });
    };

    const handleDragStart = (event: DragStartEvent) => {
        if (isSelectionMode) return;
        const { active } = event;
        const { id } = active;
        const activeContainer = findContainer(id as string);
        if (activeContainer) {
            setActiveColumn(activeContainer);
            const activeItems = columns[activeContainer];
            const card = Array.isArray(activeItems) ? activeItems.find((c) => c.id === id) : null;
            if (card) {
                setActiveCard(card);
            }
            previousColumnsRef.current = columns;
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

        queryClient.setQueryData<Record<string, CandidateCardType[]>>(["candidates", jobId], (prev) => {
            if (!prev) return prev;
            const activeItems = prev[activeContainer] || [];
            const overItems = prev[overContainer] || [];
            const activeIndex = activeItems.findIndex((item) => item.id === id);
            const overIndex = overItems.findIndex((item) => item.id === overId);

            if (activeIndex === -1) return prev;

            let newIndex;

            if (overId in prev) {
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
                    ...(prev[activeContainer] || []).filter((item: CandidateCardType) => item.id !== active.id),
                ],
                [overContainer]: [
                    ...(prev[overContainer] || []).slice(0, newIndex),
                    activeItems[activeIndex],
                    ...(prev[overContainer] || []).slice(
                        newIndex,
                        (prev[overContainer] || []).length
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
            setActiveColumn(null);
            return;
        }

        const activeContainer = findContainer(id as string);
        const overContainer = findContainer(overId as string);

        if (
            activeContainer &&
            overContainer &&
            activeContainer === overContainer
        ) {
            const activeItems = columns[activeContainer] || [];
            const activeIndex = activeItems.findIndex(
                (item) => item.id === id
            );
            const overIndex = activeItems.findIndex(
                (item) => item.id === overId
            );

            if (activeIndex !== -1 && overIndex !== -1 && activeIndex !== overIndex) {
                queryClient.setQueryData<Record<string, CandidateCardType[]>>(["candidates", jobId], (prev) => {
                    if (!prev) return prev;
                    return {
                        ...prev,
                        [activeContainer]: arrayMove(
                            prev[activeContainer] || [],
                            activeIndex,
                            overIndex
                        ),
                    };
                });

                moveCandidateApi({ 
                    candidateId: id as string, 
                    targetColumnId: activeContainer, 
                    newIndex: overIndex,
                    sourceColumnId: activeContainer 
                });
            }
        }

        if (activeColumn && overContainer && activeColumn !== overContainer) {
            const currentColumns = queryClient.getQueryData<Record<string, CandidateCardType[]>>(["candidates", jobId]);
            if (currentColumns) {
                const overItems = currentColumns[overContainer] || [];
                const finalIndex = overItems.findIndex((c) => c.id === id);
                if (finalIndex !== -1) {
                    moveCandidateApi({ 
                        candidateId: id as string, 
                        targetColumnId: overContainer, 
                        newIndex: finalIndex,
                        sourceColumnId: activeColumn as string
                    });

                    if (overContainer === "offer") {
                        fireOfferConfetti();
                    }
                }
            }
        }

        if (activeColumn && overContainer && activeColumn !== overContainer) {
            const card = (queryClient.getQueryData<Record<string, CandidateCardType[]>>(["candidates", jobId]) || {})[overContainer]?.find((c: CandidateCardType) => c.id === id);
            if (card) {
                onColumnChange?.(id as string, card, activeColumn, overContainer);
            }
        }

        setActiveCard(null);
        setActiveColumn(null);
    };

    const collisionDetectionStrategy: CollisionDetection = useCallback(
        (args) => {
            if (activeCard && activeCard.id === args.active.id) {
                return closestCorners(args);
            }

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
                            (container) => (container.id as string) in columns
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

    return {
        columns,
        activeCard,
        sensors,
        collisionDetectionStrategy,
        handleDragStart,
        handleDragOver,
        handleDragEnd,
    };
}
