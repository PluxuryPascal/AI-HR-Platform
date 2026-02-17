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
import { useGetCandidates } from "../../candidates/api/use-get-candidates";
import { useMoveCandidate } from "../../candidates/api/use-move-candidate";
import { initialColumns } from "../utils/mock-board";
import { CandidateCard as CandidateCardType, ColumnId } from "../types";
import { fireOfferConfetti } from "@/lib/confetti";

interface UseKanbanDndOptions {
    isSelectionMode: boolean;
    onColumnChange?: (candidateId: string, card: CandidateCardType, sourceColumn: ColumnId, targetColumn: ColumnId) => void;
}

export function useKanbanDnd({ isSelectionMode, onColumnChange }: UseKanbanDndOptions) {
    const queryClient = useQueryClient();
    const { data: columns = initialColumns } = useGetCandidates();
    const { mutate: moveCandidate } = useMoveCandidate();

    // Track original state for rollback
    const previousColumnsRef = useRef<Record<ColumnId, CandidateCardType[]> | null>(null);

    const [activeCard, setActiveCard] = useState<CandidateCardType | null>(null);
    const [activeColumn, setActiveColumn] = useState<ColumnId | null>(null);

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
        const activeContainer = findContainer(id as string);
        if (activeContainer) {
            setActiveColumn(activeContainer);
            const card = columns[activeContainer].find((c) => c.id === id);
            if (card) {
                setActiveCard(card);
            }
            // Snapshot current state for rollback
            if (columns) {
                previousColumnsRef.current = columns;
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

        queryClient.setQueryData<Record<ColumnId, CandidateCardType[]>>(["candidates"], (prev) => {
            if (!prev) return prev;
            const activeItems = prev[activeContainer];
            const overItems = prev[overContainer];
            const activeIndex = activeItems.findIndex((item) => item.id === id);
            const overIndex = overItems.findIndex((item) => item.id === overId);

            if (activeIndex === -1) return prev;

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
            const activeIndex = columns[activeContainer].findIndex(
                (item) => item.id === id
            );
            const overIndex = columns[overContainer].findIndex(
                (item) => item.id === overId
            );

            if (activeIndex !== overIndex) {
                queryClient.setQueryData<Record<ColumnId, CandidateCardType[]>>(["candidates"], (prev) => {
                    if (!prev) return prev;
                    return {
                        ...prev,
                        [activeContainer]: arrayMove(
                            prev[activeContainer],
                            activeIndex,
                            overIndex
                        ),
                    };
                });
            }
        }

        // Trigger mutation to persist changes
        if (activeColumn && overContainer) {
            // Get the latest state from cache to find the final index
            const currentColumns = queryClient.getQueryData<Record<ColumnId, CandidateCardType[]>>(["candidates"]);
            if (currentColumns) {
                const finalIndex = currentColumns[overContainer].findIndex((c) => c.id === id);
                if (finalIndex !== -1) {
                    moveCandidate({
                        candidateId: id as string,
                        sourceColumnId: activeColumn,
                        targetColumnId: overContainer,
                        newIndex: finalIndex,
                        optimisticSnapshot: previousColumnsRef.current || undefined,
                    });

                    // Fire confetti if moved to "offer" column
                    if (overContainer === "offer" && activeColumn !== "offer") {
                        fireOfferConfetti();
                    }
                }
            }
        }

        // Notify parent about column change for outreach
        if (activeColumn && overContainer && activeColumn !== overContainer) {
            const card = columns[overContainer].find((c) => c.id === id);
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
