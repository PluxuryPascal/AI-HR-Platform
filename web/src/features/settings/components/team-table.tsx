"use client";

import { useTranslations } from "next-intl";
import { MoreHorizontal } from "lucide-react";

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

type Role = "owner" | "recruiter" | "manager";
type Status = "active" | "pending";

export interface Member {
    id: string;
    name: string;
    email: string;
    role: Role;
    status: Status;
    avatar?: string;
    assignedJobs?: string[];
}

const mockJobs = [
    { id: "1", title: "Senior Frontend Developer" },
    { id: "2", title: "Product Manager" },
    { id: "3", title: "UI Designer" },
    { id: "4", title: "Backend Engineer" },
    { id: "5", title: "QA Specialist" },
];

const getRoleBadgeColor = (role: Role) => {
    switch (role) {
        case "owner":
            return "default";
        case "recruiter":
            return "secondary";
        case "manager":
            return "outline";
        default:
            return "default";
    }
};

interface TeamTableProps {
    members: Member[];
}

export function TeamTable({ members }: TeamTableProps) {
    const t = useTranslations("Team");

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>{t("table.name")}</TableHead>
                    <TableHead>{t("table.role")}</TableHead>
                    <TableHead>Jobs</TableHead>
                    <TableHead>{t("table.status")}</TableHead>
                    <TableHead className="text-right">{t("table.actions")}</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {members.map((member) => (
                    <TableRow key={member.id}>
                        <TableCell className="font-medium">
                            <div className="flex items-center gap-3">
                                <Avatar className="h-9 w-9">
                                    <AvatarImage src={member.avatar} alt={member.name} />
                                    <AvatarFallback>{member.name.charAt(0).toUpperCase()}</AvatarFallback>
                                </Avatar>
                                <div className="flex flex-col">
                                    <span>{member.name}</span>
                                    <span className="text-xs text-muted-foreground">{member.email}</span>
                                </div>
                            </div>
                        </TableCell>
                        <TableCell>
                            <Badge variant={getRoleBadgeColor(member.role)}>
                                {t(`roles.${member.role}`)}
                            </Badge>
                        </TableCell>
                        <TableCell>
                            {member.role === "manager" ? (
                                <div className="flex flex-wrap gap-1">
                                    {member.assignedJobs?.length ? (
                                        member.assignedJobs.map(jobId => {
                                            const job = mockJobs.find(j => j.id === jobId);
                                            return job ? (
                                                <Badge key={jobId} variant="outline" className="text-xs">
                                                    {job.title}
                                                </Badge>
                                            ) : null;
                                        })
                                    ) : (
                                        <span className="text-muted-foreground text-sm">-</span>
                                    )}
                                </div>
                            ) : (
                                <span className="text-muted-foreground text-sm">All</span>
                            )}
                        </TableCell>
                        <TableCell>
                            <Badge variant={member.status === "active" ? "default" : "secondary"} className={member.status === "active" ? "bg-green-500 hover:bg-green-600" : "bg-yellow-500 hover:bg-yellow-600"}>
                                {t(`status.${member.status}`)}
                            </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" className="h-8 w-8 p-0">
                                        <span className="sr-only">Open menu</span>
                                        <MoreHorizontal className="h-4 w-4" />
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                    <DropdownMenuItem onClick={() => navigator.clipboard.writeText(member.email)}>
                                        Copy Email
                                    </DropdownMenuItem>
                                    <DropdownMenuSeparator />
                                    <DropdownMenuItem className="text-red-600">Remove Member</DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}
