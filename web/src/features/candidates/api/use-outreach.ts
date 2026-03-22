import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";

interface GenerateEmailRequest {
    candidateId: string;
    type: "rejection" | "interview_invite";
    tone: "professional" | "friendly" | "brief";
    locale: string;
}

interface GenerateEmailResponse {
    id: string;
    candidate_id: string;
    subject: string;
    body: string;
}

interface SendEmailRequest {
    communicationId: string;
    subject: string;
    body: string;
}

export const useGenerateEmail = () => {
    return useMutation({
        mutationFn: async (req: GenerateEmailRequest) => {
            const response = await apiClient.post<GenerateEmailResponse>(
                `/outreach/${req.candidateId}/generate`,
                {
                    type: req.type,
                    tone: req.tone,
                    locale: req.locale,
                }
            );
            return response;
        },
        onError: (error: any) => {
            toast.error(error.message || "Failed to generate email");
        },
    });
};

export const useSendEmail = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (req: SendEmailRequest) => {
            const response = await apiClient.post<{ success: boolean }>(
                `/outreach/${req.communicationId}/send`,
                {
                    subject: req.subject,
                    body: req.body,
                }
            );
            return response;
        },
        onSuccess: () => {
            toast.success("Email sent successfully");
            // Optionally invalidate some queries if we have an outreach history
            // queryClient.invalidateQueries({ queryKey: ["outreach", "history"] });
        },
        onError: (error: any) => {
            toast.error(error.message || "Failed to send email");
        },
    });
};
