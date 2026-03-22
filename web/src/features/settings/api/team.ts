import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { TeamMember, TeamInvite } from "../types/team";
import { toast } from "sonner";
import { ApiResponse } from "@/types";

export function useGetTeamMembers() {
    return useQuery<TeamMember[]>({
        queryKey: ["team-members"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<TeamMember[]>>("/auth/members");
            return response.data;
        },
    });
}

export function useGetTeamInvites() {
    return useQuery<TeamInvite[]>({
        queryKey: ["team-invites"],
        queryFn: async () => {
            const response = await apiClient.get<ApiResponse<TeamInvite[]>>("/invite/list");
            return response.data;
        },
    });
}

export interface InviteUserDTO {
    email: string;
    role: string;
    job_ids?: string[];
    locale?: string;
}

export function useInviteUser() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: InviteUserDTO) => {
            const response = await apiClient.post<ApiResponse<void>>("/invite/invite", data);
            return response?.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["team-invites"] });
            toast.success("Приглашение отправлено");
        },
        onError: (error: any) => {
            console.error("Invite user error:", error);
            toast.error(error.message || "Не удалось отправить приглашение");
        },
    });
}
