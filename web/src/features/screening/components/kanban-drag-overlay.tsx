"use client";

import { DragOverlay, DropAnimation, defaultDropAnimationSideEffects } from "@dnd-kit/core";
import { CandidateCard } from "./candidate-card";
import { CandidateCard as CandidateCardType } from "../types";

const dropAnimation: DropAnimation = {
    sideEffects: defaultDropAnimationSideEffects({
        styles: {
            active: {
                opacity: "0.5",
            },
        },
    }),
};

interface KanbanDragOverlayProps {
    activeCard: CandidateCardType | null;
}

export function KanbanDragOverlay({ activeCard }: KanbanDragOverlayProps) {
    return (
        <DragOverlay dropAnimation={dropAnimation}>
            {activeCard ? <CandidateCard candidate={activeCard} isOverlay /> : null}
        </DragOverlay>
    );
}
