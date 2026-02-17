"use client";

import { useTranslations } from "next-intl";
import { Link, useRouter } from "@/i18n/routing";
import { MoreHorizontal, FileText, Mail, User } from "lucide-react";

import { Checkbox } from "@/components/ui/checkbox";
import {
    TableCell,
    TableRow,
} from "@/components/ui/table";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { MatchScoreBadge } from "./match-score-badge";
import { Candidate } from "../hooks/use-candidates";

interface CandidateRowProps {
    candidate: Candidate;
    isSelected: boolean;
    onToggleSelection: (id: string) => void;
}

export function CandidateRow({ candidate, isSelected, onToggleSelection }: CandidateRowProps) {
    const router = useRouter();

    return (
        <TableRow data-state={isSelected ? "selected" : undefined}>
            <TableCell>
                <Checkbox
                    checked={isSelected}
                    onCheckedChange={() => onToggleSelection(candidate.id)}
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
    );
}
