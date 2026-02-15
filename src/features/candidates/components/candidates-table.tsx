"use client";

import { useState } from "react";
import { Checkbox } from "@/components/ui/checkbox";
import { OutreachDrawer } from "@/features/screening/components/outreach-drawer";

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
import { CandidateScoreBadge } from "./candidate-score-badge";

export function CandidatesTable() {
    const t = useTranslations("Candidates.table");
    const tEmpty = useTranslations("Candidates.empty");
    const { data: candidates, isLoading } = useCandidates();
    const router = useRouter();

    const [selectedCandidates, setSelectedCandidates] = useState<Set<string>>(new Set());
    const [isOutreachOpen, setIsOutreachOpen] = useState(false);

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
        setIsOutreachOpen(true);
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
        <div className="space-y-4">
            {selectedCandidates.size > 0 && (
                <div className="bg-primary/5 border border-primary/20 rounded-lg p-4 flex items-center justify-between animate-in fade-in slide-in-from-top-2">
                    <div className="flex items-center gap-2">
                        <Badge variant="default" className="bg-primary text-primary-foreground">
                            {selectedCandidates.size}
                        </Badge>
                        <span className="text-sm font-medium text-primary">Selected</span>
                    </div>
                    <div className="flex items-center gap-2">
                        <Button size="sm" variant="destructive" onClick={handleBulkReject}>
                            {t("action.bulkReject", { count: selectedCandidates.size })}
                        </Button>
                        <Button size="sm" variant="ghost" onClick={() => setSelectedCandidates(new Set())}>
                            Cancel
                        </Button>
                    </div>
                </div>
            )}

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
                                <TableCell>{candidate.appliedRole}</TableCell>
                                <TableCell>
                                    <div className="flex flex-col gap-1">
                                        <CandidateScoreBadge score={candidate.score} />
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

            <OutreachDrawer
                isOpen={isOutreachOpen}
                onClose={() => setIsOutreachOpen(false)}
                candidate={selectedCandidatesList.map(c => ({
                    id: c.id,
                    name: c.name,
                    role: c.appliedRole,
                    score: c.score,
                    // Avatar URL not in Candidate type
                }))}
                type="rejection"
            />
        </div>
    );
}
