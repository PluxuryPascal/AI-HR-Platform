export interface Stage {
    id: string;
    job_id?: string;
    team_id: string;
    code: string;
    title: string;
    position: number;
    is_terminal: boolean;
    color?: string;
}

export interface CreateStageRequest {
    title: string;
    code: string;
    position: number;
    is_terminal?: boolean;
    color?: string;
}

export interface UpdateStageRequest {
    title?: string;
    code?: string;
    position?: number;
    is_terminal?: boolean;
    color?: string;
}
