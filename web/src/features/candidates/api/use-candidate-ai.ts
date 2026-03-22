import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { ApiResponse } from "@/types";

export interface ChatMessage {
    id?: string;
    session_id?: string;
    role: "user" | "assistant" | "system";
    content: string;
    created_at?: string;
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

export interface ChatSession {
    id: string;
    team_id: string;
    user_id: string;
    type: string;
    target_candidate_id?: string;
    title: string;
    created_at: string;
    updated_at: string;
}

// ChatMessage is defined at top of file

export function useCandidateChat() {
    return useMutation({
        mutationFn: async (payload: ChatRequest) => {
            const response = await apiClient.post<ApiResponse<ChatResponse>>("/chat", payload);
            return response.data;
        },
    });
}

export function useGetChatSessions() {
    return useQuery<ChatSession[]>({
        queryKey: ["chat", "sessions"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<ChatSession[]>>("/chat/sessions");
            return response.data;
        },
    });
}

export function useGetChatHistory(sessionId: string) {
    return useQuery<ChatMessage[]>({
        queryKey: ["chat", "history", sessionId],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<ChatMessage[]>>(`/chat/sessions/${sessionId}/history`);
            return response.data;
        },
        enabled: !!sessionId,
    });
}

export function useGenerateInterview() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, locale }: { id: string; locale: string }) => {
            const response = await apiClient.post<ApiResponse<InterviewGuide>>(`/interview/${id}/questions`, { locale });
            return response.data;
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
    return useQuery<InterviewGuide | null>({
        queryKey: ["candidate", candidateId, "interview"],
        queryFn: async () => {
            try {
                const response = await apiClient.get<ApiResponse<InterviewGuide>>(`/interview/${candidateId}/questions`);
                return response.data;
            } catch (err) {
                console.warn("Interview guide not found or error:", err);
                return null;
            }
        },
        enabled: !!candidateId,
        retry: false,
    });
}
