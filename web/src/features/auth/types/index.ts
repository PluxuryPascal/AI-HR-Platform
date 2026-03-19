export interface RegisterRequest {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    team_name: string;
}

export interface AuthResponse {
    token: string;
    user: {
        id: string;
        email: string;
        first_name: string;
        last_name: string;
        role: string;
        team_id: string;
        team_name: string;
    };
}
