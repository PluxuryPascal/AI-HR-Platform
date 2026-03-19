import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { TeamMember, TeamInvite } from "../types/team";
import { toast } from "sonner";

export function useGetTeamMembers() {
    return useQuery<TeamMember[]>({
        queryKey: ["team-members"],
        queryFn: async () => {
            return await apiClient.get<TeamMember[]>("/auth/members");
        },
    });
}

export function useGetTeamInvites() {
    return useQuery<TeamInvite[]>({
        queryKey: ["team-invites"],
        queryFn: async () => {
            return await apiClient.get<TeamInvite[]>("/invite/list");
        },
    });
}

export interface InviteUserDTO {
    email: string;
    role: string;
    job_ids?: string[];
}

export function useInviteUser() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: InviteUserDTO) => {
            return await apiClient.post("/invite/invite", data);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["team-invites"] });
            toast.success("Приглашение отправлено");
        },
        onError: () => {
            toast.error("Не удалось отправить приглашение");
        },
    });
}
