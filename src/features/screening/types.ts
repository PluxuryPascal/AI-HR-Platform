export type ColumnId = 'new' | 'screening' | 'interview' | 'offer' | 'rejected';

export interface ScoreFactor {
    id: string;
    type: "positive" | "negative" | "neutral";
    text: string;
    impact: string;
}

export type CandidateCard = {
    id: string;
    name: string;
    role: string;
    score: number;
    avatarUrl?: string; // Optional, might be generated or fetched
    email?: string;
    appliedDate?: string;
    matchSummary?: string;
    scoreBreakdown?: ScoreFactor[];
    status?: string; // stored in column key but good to have
};
