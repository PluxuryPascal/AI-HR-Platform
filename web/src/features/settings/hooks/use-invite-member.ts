import { useState, useCallback } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { inviteSchema, InviteFormValues } from "../components/invite-form";
import { Member } from "../components/team-table";

const initialMembers: Member[] = [
    {
        id: "1",
        name: "Alice Johnson",
        email: "alice@example.com",
        role: "owner",
        status: "active",
        avatar: "https://i.pravatar.cc/150?u=alice",
    },
    {
        id: "2",
        name: "Bob Smith",
        email: "bob@example.com",
        role: "recruiter",
        status: "active",
        avatar: "https://i.pravatar.cc/150?u=bob",
    },
    {
        id: "3",
        name: "Charlie Brown",
        email: "charlie@example.com",
        role: "manager",
        status: "pending",
        avatar: "https://i.pravatar.cc/150?u=charlie",
        assignedJobs: ["1", "3"],
    },
];

export function useInviteMember() {
    const [members, setMembers] = useState<Member[]>(initialMembers);
    const [isInviteOpen, setIsInviteOpen] = useState(false);

    const form = useForm<InviteFormValues>({
        resolver: zodResolver(inviteSchema),
        defaultValues: {
            email: "",
            role: "recruiter",
            jobs: [],
        },
    });

    const onSubmit = useCallback((values: InviteFormValues) => {
        const newMember: Member = {
            id: Math.random().toString(36).substr(2, 9),
            name: values.email.split("@")[0], // Placeholder name
            email: values.email,
            role: values.role,
            status: "pending",
            assignedJobs: values.role === "manager" ? values.jobs : undefined,
        };

        setMembers(prev => [...prev, newMember]);
        setIsInviteOpen(false);
        form.reset();
        return newMember;
    }, [form]);

    return {
        members,
        isInviteOpen,
        setIsInviteOpen,
        form,
        onSubmit,
    };
}
