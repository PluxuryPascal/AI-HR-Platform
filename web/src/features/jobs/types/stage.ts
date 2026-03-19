export interface Stage {
    id: string;
    job_id: string;
    team_id: string;
    title: string;
    position: number;
    is_system?: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreateStageRequest {
    title: string;
    position: number;
}
