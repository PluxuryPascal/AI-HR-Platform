export type CandidateParsingStatus = "pending" | "processing" | "needs_review" | "completed" | "failed";

export interface Candidate {
    id: string;
    job_id: string;
    first_name?: string;
    last_name?: string;
    email?: string;
    resume_file_key?: string;
    parsed_text?: string;
    location?: string;
    skills?: string[];
    parsing_status: CandidateParsingStatus;
    stage_id?: string;
    created_at: string;
    updated_at?: string;
}

export interface CandidateProfile {
    candidate_id: string;
    structured_data?: string;
    ai_parsed_at?: string;
    updated_at?: string;
    missing_fields?: string[];
}

export interface ScoreFactor {
    id: string;
    candidate_id: string;
    type: "positive" | "negative" | "neutral";
    description: string;
    impact: number;
}

export interface CandidateScore {
    candidate_id: string;
    match_score: number;
    analyzed_at?: string;
}

export interface CandidateDetailResponse {
    candidate: Candidate;
    profile: CandidateProfile;
    stage_id?: string;
    score?: CandidateScore;
    factors: ScoreFactor[];
}

export interface CandidateFilter {
    first_name?: string;
    last_name?: string;
    email?: string;
    current_stage_id?: string;
}

export interface MoveCandidateRequest {
    to_stage_id: string;
    new_position: number;
}
