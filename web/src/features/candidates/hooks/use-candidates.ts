"use client";

import { useMemo } from "react";
import { useGetCandidates } from "../api/use-get-candidates";
import { CandidateCard, ColumnId } from "@/features/screening/types";

export type Candidate = CandidateCard;

export function useCandidates() {
    const { data: columns, isLoading } = useGetCandidates();

    const candidates = useMemo(() => {
        if (!columns) return [];
        const list: Candidate[] = [];
        (Object.entries(columns) as [ColumnId, CandidateCard[]][]).forEach(([colId, cards]) => {
            cards.forEach(card => {
                // Ensure status is set based on column if missing
                list.push({
                    ...card,
                    status: card.status || colId
                });
            });
        });
        return list;
    }, [columns]);

    return { data: candidates, isLoading };
}
