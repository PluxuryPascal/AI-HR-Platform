import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types";

export interface DashboardStats {
    total_candidates: number;
    total_candidates_delta: number;
    active_jobs: number;
    active_jobs_delta: number;
    upcoming_interviews: number;
    interviews_delta: number;
    avg_time_to_hire_days: number;
    avg_time_to_hire_delta: number;
}

export interface ChartDataPoint {
    date: string;
    count: number;
}

export interface ActivityLogEntry {
    id: string;
    actor_type: string;
    actor_id: string;
    action_code: string;
    target_id: string;
    created_at: string;
    candidate_first_name?: string;
    candidate_last_name?: string;
    job_title?: string;
    actor_name: string;
    match_score?: number;
}

export function useGetDashboardStats() {
    return useQuery<DashboardStats>({
        queryKey: ["dashboard", "stats"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<DashboardStats>>("/dashboard/stats");
            return response.data;
        },
    });
}

export function useGetDashboardDynamics(startDate: string, endDate: string, type: "daily" | "monthly") {
    return useQuery<ChartDataPoint[]>({
        queryKey: ["dashboard", "dynamics", startDate, endDate, type],
        queryFn: async () => {
            const params = new URLSearchParams({
                start_date: startDate,
                end_date: endDate,
                type: type,
            });
            const response = await apiClient.get<ApiResponse<ChartDataPoint[]>>(`/dashboard/applications-chart?${params.toString()}`);
            return response.data;
        },
    });
}

export function useGetRecentActivity(limit: number = 10) {
    return useQuery<ActivityLogEntry[]>({
        queryKey: ["dashboard", "activity", limit],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<ActivityLogEntry[]>>(`/dashboard/recent-activity?limit=${limit}`);
            return response.data;
        },
    });
}
