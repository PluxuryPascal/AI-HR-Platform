export interface Department {
    id: string;
    name: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateDepartmentPayload {
    name: string;
    description?: string;
}
