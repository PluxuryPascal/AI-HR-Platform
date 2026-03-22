import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types";

export interface AISettings {
    team_id: string;
    api_key?: string;
    parse_model?: string;
    score_model?: string;
    embed_model?: string;
    chat_model?: string;
}

export interface AIModel {
    id: string;
    name: string;
    provider?: string;
    is_embedding?: boolean;
}

export function useGetAISettings() {
    return useQuery({
        queryKey: ["ai-settings"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<AISettings>>("/ai-settings");
            return response.data;
        },
    });
}

export function useGetAIModels() {
    return useQuery({
        queryKey: ["ai-models"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<AIModel[]>>("/ai-settings/models");
            return response.data;
        },
    });
}

export function useUpdateAISetting() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ field, value }: { field: string; value: string }) => {
            // field can be api_key, parse_model, score_model, embed_model, chat_model
            const endpointField = field.replace(/_/g, "-");
            const endpoint = `/ai-settings/${endpointField}`;
            
            const body = field === "api_key" ? { api_key: value } : { value };
            
            const response = await apiClient.patch<ApiResponse<{ status: string }>>(endpoint, body);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["ai-settings"] });
        },
    });
}
