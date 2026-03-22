import { useState, useCallback, useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useLocale } from "next-intl";
import { inviteSchema, InviteFormValues } from "../components/invite-form";
import { Member, Role } from "../components/team-table";
import { useGetTeamMembers, useGetTeamInvites, useInviteUser } from "../api/team";

export function useInviteMember() {
    const [isInviteOpen, setIsInviteOpen] = useState(false);
    const locale = useLocale();
    
    const { data: membersRaw, isLoading: isLoadingMembers } = useGetTeamMembers();
    const { data: invitesRaw, isLoading: isLoadingInvites } = useGetTeamInvites();
    const inviteMutation = useInviteUser();

    const members = useMemo(() => {
        const result: Member[] = [];
        const seenIds = new Set<string>();

        if (membersRaw) {
            membersRaw.forEach(m => {
                if (!m.id || seenIds.has(m.id)) return;
                seenIds.add(m.id);
                result.push({
                    id: m.id,
                    name: `${m.first_name || ""} ${m.last_name || ""}`.trim() || m.email?.split('@')[0] || "Unknown",
                    email: m.email || "",
                    role: (m.role as Role) || "recruiter",
                    status: "active",
                    avatar: m.avatar || `https://i.pravatar.cc/150?u=${m.id}`,
                });
            });
        }
        
        if (invitesRaw) {
            invitesRaw.forEach(i => {
                if (!i.id || seenIds.has(i.id)) return;
                seenIds.add(i.id);
                result.push({
                    id: i.id,
                    name: i.email?.split('@')[0] || "Unknown",
                    email: i.email || "",
                    role: (i.role as Role) || "recruiter",
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
        try {
            await inviteMutation.mutateAsync({
                email: values.email,
                role: values.role,
                job_ids: values.jobs && values.jobs.length > 0 ? values.jobs : undefined,
                locale: locale,
            });
            
            setIsInviteOpen(false);
            form.reset();
        } catch (error: any) {
            // Error is already shown in mutation onError but we catch it here 
            // to prevent bubbling up and causing runtime error overlay
            console.error("Invite member error:", error);
        }
    }, [inviteMutation, form, locale]);

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
