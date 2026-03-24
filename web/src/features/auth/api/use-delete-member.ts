import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";

export function useDeleteMember() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (id: string) => {
            return await apiClient.delete(`/members/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["members"] });
            toast.success("Member removed from team");
        },
        onError: (error: any) => {
            toast.error(error?.response?.data?.message || "Failed to remove member");
        }
    });
}
