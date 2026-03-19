import { useState, useCallback, useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { inviteSchema, InviteFormValues } from "../components/invite-form";
import { Member, Role } from "../components/team-table";
import { useGetTeamMembers, useGetTeamInvites, useInviteUser } from "../api/team";

export function useInviteMember() {
    const [isInviteOpen, setIsInviteOpen] = useState(false);
    
    const { data: membersRaw, isLoading: isLoadingMembers } = useGetTeamMembers();
    const { data: invitesRaw, isLoading: isLoadingInvites } = useGetTeamInvites();
    const inviteMutation = useInviteUser();

    const members = useMemo(() => {
        const result: Member[] = [];
        
        if (membersRaw) {
            membersRaw.forEach(m => {
                result.push({
                    id: m.id,
                    name: `${m.first_name || ""} ${m.last_name || ""}`.trim() || m.email.split('@')[0],
                    email: m.email,
                    role: m.role as Role,
                    status: "active",
                    avatar: m.avatar || `https://i.pravatar.cc/150?u=${m.id}`,
                });
            });
        }
        
        if (invitesRaw) {
            invitesRaw.forEach(i => {
                result.push({
                    id: i.id,
                    name: i.email.split('@')[0],
                    email: i.email,
                    role: i.role as Role,
                    status: "pending",
                });
            });
        }
        
        return result;
    }, [membersRaw, invitesRaw]);

    const form = useForm<InviteFormValues>({
        resolver: zodResolver(inviteSchema),
        defaultValues: {
            email: "",
            role: "recruiter",
            jobs: [],
        },
    });

    const onSubmit = useCallback(async (values: InviteFormValues) => {
        await inviteMutation.mutateAsync({
            email: values.email,
            role: values.role,
            job_ids: values.jobs && values.jobs.length > 0 ? values.jobs : undefined,
        });
        
        setIsInviteOpen(false);
        form.reset();
    }, [inviteMutation, form]);

    return {
        members,
        isLoading: isLoadingMembers || isLoadingInvites,
        isInviteOpen,
        setIsInviteOpen,
        form,
        onSubmit,
        isSubmitting: inviteMutation.isPending,
    };
}
