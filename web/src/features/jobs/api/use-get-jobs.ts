import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Job, JobFilter } from "../types/job";
import { PaginatedResponse } from "@/types";

interface GetJobsParams {
    page?: number;
    per_page?: number;
    filter?: JobFilter;
}

export function useGetJobs({ page = 1, per_page = 10, filter }: GetJobsParams = {}) {
    return useQuery({
        queryKey: ["jobs", { page, per_page, filter }],
        queryFn: async () => {
            const response = await apiClient.post<PaginatedResponse<Job>>("/jobs/list", {
                pagination: { page, per_page },
                filter,
            });
            return response;
        },
        refetchInterval: 5000,
    });
}
