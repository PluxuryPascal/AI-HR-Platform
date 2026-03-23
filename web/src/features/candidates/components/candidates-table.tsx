"use client";
import { useState, useMemo, useEffect } from "react";

import { Checkbox } from "@/components/ui/checkbox";
import { useTranslations } from "next-intl";
import { Users } from "lucide-react";

import {
    Table,
    TableBody,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { UploadResumeDialog } from "./upload-resume-dialog";

import { useCandidatesPagination } from "../hooks/use-candidates";
import { BulkActionBar } from "./bulk-action-bar";
import { CandidateRow } from "./candidate-row";
import { useTableSelection } from "../hooks/use-table-selection";
import { useBulkActions } from "../hooks/use-bulk-actions";
import { useGetStages } from "@/features/jobs/api/use-get-stages";
import { Pagination } from "@/components/shared/pagination";
import { ChevronUp, ChevronDown, ChevronsUpDown } from "lucide-react";
import { CandidateFilter } from "../types/candidate";

import { OutreachDrawer } from "@/features/screening/components/outreach-drawer";

export function CandidatesTable({ jobId, filter }: { jobId: string; filter?: CandidateFilter }) {
    const t = useTranslations("Candidates.table");
    const tEmpty = useTranslations("Candidates.empty");
    
    const [page, setPage] = useState(1);
    const [perPage] = useState(10);
    const [sortId, setSortId] = useState<string>("created_at");
    const [sortDesc, setSortDesc] = useState(true);

    // Reset to page 1 when filter or job changes
    useEffect(() => {
        setPage(1);
    }, [jobId, JSON.stringify(filter)]);

    const { data: candidates, meta, isLoading } = useCandidatesPagination({ 
        jobId, 
        page, 
        per_page: perPage,
        filter: {
            ...filter,
            sort: { sort_id: sortId, sort_desc: sortDesc }
        }
    });

    const { data: stages } = useGetStages(jobId);

    const {
        selectedIds,
        toggleSelection,
        toggleAll,
        clearSelection,
        isAllSelected,
        selectedItems,
    } = useTableSelection(candidates);

    const handleSort = (id: string) => {
        if (sortId === id) {
            setSortDesc(!sortDesc);
        } else {
            setSortId(id);
            setSortDesc(true);
        }
        setPage(1);
    };

    const SortIcon = ({ id }: { id: string }) => {
        if (sortId !== id) return <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />;
        return sortDesc ? <ChevronDown className="ml-2 h-4 w-4" /> : <ChevronUp className="ml-2 h-4 w-4" />;
    };

    const {
        bulkActionType,
        isOutreachOpen,
        handleBulkReject,
        handleBulkMove,
        handleSendEmail,
        closeOutreach,
    } = useBulkActions({ selectedIds, clearSelection, jobId, stages });

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
                    <UploadResumeDialog jobId={jobId}>
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
                                    checked={isAllSelected}
                                    onCheckedChange={toggleAll}
                                />
                            </TableHead>
                            <TableHead className="cursor-pointer hover:bg-muted/50 transition-colors" onClick={() => handleSort("first_name")}>
                                <div className="flex items-center">
                                    {t("name")}
                                    <SortIcon id="first_name" />
                                </div>
                            </TableHead>
                            <TableHead>{t("role")}</TableHead>
                            <TableHead className="cursor-pointer hover:bg-muted/50 transition-colors" onClick={() => handleSort("match_score")}>
                                <div className="flex items-center">
                                    {t("score")}
                                    <SortIcon id="match_score" />
                                </div>
                            </TableHead>
                            <TableHead>{t("status")}</TableHead>
                            <TableHead className="text-right"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {candidates.map((candidate) => (
                            <CandidateRow
                                key={candidate.id}
                                candidate={candidate}
                                isSelected={selectedIds.has(candidate.id)}
                                onToggleSelection={toggleSelection}
                            />
                        ))}
                    </TableBody>
                </Table>
                <Pagination 
                    currentPage={page} 
                    totalPages={meta?.total_pages || 1} 
                    onPageChange={setPage} 
                />
            </div>

            <BulkActionBar
                selectedCount={selectedIds.size}
                onClear={clearSelection}
                onRejectAll={handleBulkReject}
                onMoveTo={handleBulkMove}
                stages={stages}
            />

            <OutreachDrawer
                isOpen={isOutreachOpen}
                onClose={closeOutreach}
                candidate={selectedItems}
                type={bulkActionType || "rejection"}
                onSend={handleSendEmail}
            />
        </div>
    );
}
