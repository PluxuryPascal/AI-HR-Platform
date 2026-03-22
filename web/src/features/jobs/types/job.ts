export type JobStatus = "status_draft" | "status_published" | "status_closed" | "status_archived";

export type WorkFormat = "remote" | "office" | "hybrid";

export interface Job {
    id: string;
    team_id: string;
    title: string;
    department_id?: string;
    department_name?: string;
    work_format: WorkFormat;
    description?: string;
    requirements?: string[];
    status: JobStatus;
    salary_min?: number;
    salary_max?: number;
    currency: string;
    created_at: string;
    updated_at: string;
}

export interface JobFilter {
    title?: string;
    department_name?: string;
    status?: JobStatus;
    work_format?: WorkFormat;
}

export interface CreateJobRequest {
    title: string;
    department_id?: string;
    work_format: WorkFormat;
    description?: string;
    requirements?: string[];
    salary_min?: number;
    salary_max?: number;
    currency?: string;
}

export interface UpdateJobRequest extends Partial<CreateJobRequest> {}
