export type ColumnId = 'new' | 'screening' | 'interview' | 'offer' | 'rejected';

export type CandidateCard = {
    id: string;
    name: string;
    role: string;
    score: number;
    avatarUrl?: string; // Optional, might be generated or fetched
};
