export interface TeamMember {
    id: string;
    team_id: string;
    email: string;
    first_name: string;
    last_name: string;
    role: string;
    created_at: string;
    updated_at: string;
    avatar?: string;
}

export interface TeamInvite {
    id: string;
    team_id: string;
    email: string;
    role: string;
    token: string;
    status: "pending" | "accepted" | "expired" | "failed";
    expires_at: string;
    created_at: string;
}
