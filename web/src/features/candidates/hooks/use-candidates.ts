import { useMemo } from "react";
import { useGetCandidates } from "../api/use-get-candidates";
import { CandidateCard, ColumnId } from "@/features/screening/types";
import { Candidate as BackendCandidate } from "../types/candidate";

export type Candidate = CandidateCard;

export function useCandidates(jobId?: string) {
    const { data: columns, isLoading } = useGetCandidates(jobId);

    const candidates = useMemo(() => {
        if (!columns) return [];
        const list: Candidate[] = [];
        (Object.entries(columns) as [string, BackendCandidate[]][]).forEach(([colId, cards]) => {
            cards.forEach(card => {
                list.push({
                    id: card.id,
                    name: `${card.first_name || ""} ${card.last_name || ""}`.trim() || "Anonymous",
                    role: "Candidate",
                    score: 0,
                    email: card.email,
                    status: card.stage_id || colId
                });
            });
        });
        return list;
    }, [columns]);

    return { data: candidates, isLoading };
}
