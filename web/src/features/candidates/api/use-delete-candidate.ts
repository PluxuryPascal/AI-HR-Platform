import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

import { apiClient } from "@/lib/api-client";

export const deleteCandidate = (id: string): Promise<void> => {
    return apiClient.delete(`/candidates/${id}`);
};

export const useDeleteCandidate = (jobId?: string) => {
    const queryClient = useQueryClient();
    const t = useTranslations("Candidates.mutations");

    return useMutation({
        mutationFn: deleteCandidate,
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: jobId 
                    ? ["candidates", jobId] 
                    : ["candidates"],
            });
            toast.success(t("deleteSuccess"));
        },
        onError: (error: any) => {
            toast.error(error.message || t("deleteError"));
        },
    });
};
