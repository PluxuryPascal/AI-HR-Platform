import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Job } from "../types/job";

export function useGetJob(id: string) {
    return useQuery({
        queryKey: ["job", id],
        queryFn: async () => {
            const response = await apiClient.get<{ data: Job; success: boolean }>(`/jobs/${id}`);
            
            if (!response?.success || !response?.data) {
                throw new Error("Failed to fetch job data");
            }

            const job = response.data as any;
            
            // Backend returns requirements as extracted_requirements (JSON array)
            if (job.extracted_requirements && !job.requirements) {
                if (Array.isArray(job.extracted_requirements)) {
                    job.requirements = job.extracted_requirements;
                } else if (typeof job.extracted_requirements === 'string') {
                    try {
                        job.requirements = JSON.parse(job.extracted_requirements);
                    } catch (e) {
                        console.warn("Failed to parse requirements string:", e);
                        job.requirements = [];
                    }
                } else {
                    // Handle other cases like base64 bytes if needed, 
                    // but RawMessage should be an array/object/string.
                    job.requirements = [];
                }
            }
            
            return job as Job;
        },
        enabled: !!id,
    });
}
