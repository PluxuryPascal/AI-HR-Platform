import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

import { apiClient } from "@/lib/api-client";

interface BulkDeleteCandidatesOptions {
    candidateIds: string[];
}

export const bulkDeleteCandidates = ({
    jobId,
    data,
}: {
    jobId?: string;
    data: BulkDeleteCandidatesOptions;
}): Promise<void> => {
    const url = jobId 
        ? `/jobs/${jobId}/candidates/bulk-delete` 
        : `/candidates/bulk-delete`;
        
    return apiClient.post(url, {
        candidate_ids: data.candidateIds,
    });
};

export const useBulkDeleteCandidates = (jobId?: string) => {
    const queryClient = useQueryClient();
    const t = useTranslations("Candidates.mutations");

    return useMutation({
        mutationFn: (data: BulkDeleteCandidatesOptions) => bulkDeleteCandidates({ jobId, data }),
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: jobId 
                    ? ["candidates", jobId] 
                    : ["candidates"],
            });
            toast.success(t("bulkDeleteSuccess"));
        },
        onError: (error: any) => {
            toast.error(error.message || t("bulkDeleteError"));
        },
    });
};
