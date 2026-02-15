"use client";

import { useState } from "react";
import { Checkbox } from "@/components/ui/checkbox";

import { useTranslations } from "next-intl";
import { Link, useRouter } from "@/i18n/routing";
import { MoreHorizontal, FileText, Mail, User, Users } from "lucide-react";

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
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { EmptyState } from "@/components/shared/empty-state";
import { UploadResumeDialog } from "./upload-resume-dialog";

import { useCandidates, Candidate } from "../hooks/use-candidates";
import { useBulkUpdateCandidates } from "../api/use-bulk-update-candidates";
import { BulkActionBar } from "./bulk-action-bar";
import { MatchScoreBadge } from "./match-score-badge";
import { ColumnId } from "@/features/screening/types";

import { OutreachDrawer } from "@/features/screening/components/outreach-drawer";

export function CandidatesTable() {
    const t = useTranslations("Candidates.table");
    const tEmpty = useTranslations("Candidates.empty");
    const { data: candidates, isLoading } = useCandidates();
    const router = useRouter();

    const [selectedCandidates, setSelectedCandidates] = useState<Set<string>>(new Set());
    const [bulkActionType, setBulkActionType] = useState<"rejection" | "invitation" | null>(null);

    const { mutate: bulkUpdate } = useBulkUpdateCandidates();

    const toggleSelection = (id: string) => {
        const newSelection = new Set(selectedCandidates);
        if (newSelection.has(id)) {
            newSelection.delete(id);
        } else {
            newSelection.add(id);
        }
        setSelectedCandidates(newSelection);
    };

    const toggleAll = () => {
        if (selectedCandidates.size === candidates.length) {
            setSelectedCandidates(new Set());
        } else {
            setSelectedCandidates(new Set(candidates.map(c => c.id)));
        }
    };

    const handleBulkReject = () => {
        setBulkActionType("rejection");
    };

    const handleBulkMove = (status: string) => {
        if (status === "interview") {
            setBulkActionType("invitation");
            return;
        }

        bulkUpdate({
            candidateIds: Array.from(selectedCandidates),
            newStatus: status as ColumnId
        });
        setSelectedCandidates(new Set());
    };

    const handleSendEmail = (content: string) => {
        // In a real app, we would send the email here using the 'content'
        // For now, we just update the status
        const newStatus = bulkActionType === "rejection" ? "rejected" : "interview";

        bulkUpdate({
            candidateIds: Array.from(selectedCandidates),
            newStatus: newStatus
        });
        setSelectedCandidates(new Set());
        setBulkActionType(null);
    };

    const selectedCandidatesList = candidates.filter(c => selectedCandidates.has(c.id));

    if (isLoading) {
        return (
            <div className="space-y-3">
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
            </div>
        );
    }


    if (candidates.length === 0) {
        return (
            <EmptyState
                icon={Users}
                title={tEmpty("title")}
                description={tEmpty("desc")}
                action={
                    <UploadResumeDialog>
                        <Button>
                            {tEmpty("button")}
                        </Button>
                    </UploadResumeDialog>
                }
                className="border rounded-md"
            />
        );
    }

    return (
        <div className="space-y-4 pb-20"> {/* Added padding bottom for FAB */}
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[50px]">
                                <Checkbox
                                    checked={selectedCandidates.size === candidates.length && candidates.length > 0}
                                    onCheckedChange={toggleAll}
                                />
                            </TableHead>
                            <TableHead>{t("name")}</TableHead>
                            <TableHead>{t("role")}</TableHead>
                            <TableHead>{t("score")}</TableHead>
                            <TableHead>{t("status")}</TableHead>
                            <TableHead className="text-right"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {candidates.map((candidate) => (
                            <TableRow key={candidate.id} data-state={selectedCandidates.has(candidate.id) ? "selected" : undefined}>
                                <TableCell>
                                    <Checkbox
                                        checked={selectedCandidates.has(candidate.id)}
                                        onCheckedChange={() => toggleSelection(candidate.id)}
                                    />
                                </TableCell>
                                <TableCell className="font-medium">
                                    <Link href={`/dashboard/candidates/${candidate.id}`} className="flex items-center gap-3 hover:underline">
                                        <Avatar className="h-9 w-9">
                                            <AvatarImage src={`https://api.dicebear.com/7.x/initials/svg?seed=${candidate.name}`} />
                                            <AvatarFallback>{candidate.name.substring(0, 2).toUpperCase()}</AvatarFallback>
                                        </Avatar>
                                        <div className="flex flex-col">
                                            <span>{candidate.name}</span>
                                            <span className="text-xs text-muted-foreground">{candidate.email}</span>
                                        </div>
                                    </Link>
                                </TableCell>
                                <TableCell>{candidate.role}</TableCell>
                                <TableCell>
                                    <div className="flex flex-col gap-1">
                                        <MatchScoreBadge score={candidate.score} breakdown={candidate.scoreBreakdown} candidateId={candidate.id} />
                                        <span className="text-xs text-muted-foreground max-w-[200px] truncate" title={candidate.matchSummary}>
                                            {candidate.matchSummary}
                                        </span>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="secondary" className="capitalize">
                                        {candidate.status}
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
                                            <DropdownMenuItem onClick={() => router.push(`/dashboard/candidates/${candidate.id}`)}>
                                                <User className="mr-2 h-4 w-4" /> View Details
                                            </DropdownMenuItem>
                                            <DropdownMenuItem>
                                                <FileText className="mr-2 h-4 w-4" /> View Resume
                                            </DropdownMenuItem>
                                            <DropdownMenuItem>
                                                <Mail className="mr-2 h-4 w-4" /> Email Candidate
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>

            <BulkActionBar
                selectedCount={selectedCandidates.size}
                onClear={() => setSelectedCandidates(new Set())}
                onRejectAll={handleBulkReject}
                onMoveTo={handleBulkMove}
            />

            <OutreachDrawer
                isOpen={!!bulkActionType}
                onClose={() => setBulkActionType(null)}
                candidate={selectedCandidatesList}
                type={bulkActionType || "rejection"}
                onSend={handleSendEmail}
            />
        </div>
    );
}
