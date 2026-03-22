import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { Candidate } from "../types/candidate";
import { ApiResponse } from "@/types";

interface UploadResumeParams {
    jobId: string;
    file: File;
    locale?: string;
}

export function useUploadResume() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ jobId, file, locale = "ru" }: UploadResumeParams) => {
            const formData = new FormData();
            formData.append("file", file);
            formData.append("locale", locale);

            // Fetch doesn't automatically set Content-Type correctly for FormData if we set it manually.
            // Our apiClient setting Content-Type to application/json by default might break this.
            // Let's check apiClient implementation... 
            // It has headers: { "Content-Type": "application/json", ...headers }.
            // We need to override it.
            
            const response = await apiClient.post<ApiResponse<{ candidate: Candidate }>>(
                `/jobs/${jobId}/candidates/upload`,
                formData,
                {
                    headers: {
                        // Let the browser set the boundary
                        "Content-Type": "", 
                    }
                } as any
            );
            
            return response.data;
        },
        onSuccess: (_, { jobId }) => {
            queryClient.invalidateQueries({ queryKey: ["candidates", jobId] });
        },
        onError: (error: any) => {
            console.error("Upload resume error:", error);
            toast.error(error.message || "Ошибка при загрузке резюме");
        },
    });
}
