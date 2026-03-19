import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";

export interface ChatMessage {
    role: "user" | "assistant" | "system";
    content: string;
}

export interface ChatRequest {
    candidate_id?: string;
    question: string;
    locale: string;
}

export interface ChatResponse {
    answer: string;
    session_id?: string;
}

export interface InterviewQuestion {
    id: string;
    question: string;
    answer: string;
    category: string;
}

export interface InterviewGuide {
    id: string;
    candidate_id: string;
    job_id: string;
    questions: InterviewQuestion[];
    created_at: string;
}

export function useCandidateChat() {
    return useMutation({
        mutationFn: async (payload: ChatRequest) => {
            const response = await apiClient.post<ChatResponse>("/chat", payload);
            return response;
        },
    });
}

export function useGenerateInterview() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, locale }: { id: string; locale: string }) => {
            const response = await apiClient.post<InterviewGuide>(`/candidates/${id}/interview/questions`, { locale });
            return response;
        },
        onSuccess: (data, variables) => {
            queryClient.invalidateQueries({ queryKey: ["candidate", variables.id, "interview"] });
            toast.success("Гайд для интервью сгенерирован");
        },
        onError: () => {
            toast.error("Ошибка при генерации интервью");
        },
    });
}

export function useGetInterviewGuide(candidateId: string) {
    return useQuery<InterviewGuide>({
        queryKey: ["candidate", candidateId, "interview"],
        queryFn: async () => {
            // This might be a GET or the result is cached from previous POST
            // For now, let's assume it's part of candidate details or separate GET
            const response = await apiClient.get<InterviewGuide>(`/candidates/${candidateId}/interview/questions`);
            return response;
        },
        enabled: !!candidateId,
    });
}
