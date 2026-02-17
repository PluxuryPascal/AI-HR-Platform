"use client";

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

import { useCandidates } from "../hooks/use-candidates";
import { BulkActionBar } from "./bulk-action-bar";
import { CandidateRow } from "./candidate-row";
import { useTableSelection } from "../hooks/use-table-selection";
import { useBulkActions } from "../hooks/use-bulk-actions";

import { OutreachDrawer } from "@/features/screening/components/outreach-drawer";

export function CandidatesTable() {
    const t = useTranslations("Candidates.table");
    const tEmpty = useTranslations("Candidates.empty");
    const { data: candidates, isLoading } = useCandidates();

    const {
        selectedIds,
        toggleSelection,
        toggleAll,
        clearSelection,
        isAllSelected,
        selectedItems,
    } = useTableSelection(candidates);

    const {
        bulkActionType,
        isOutreachOpen,
        handleBulkReject,
        handleBulkMove,
        handleSendEmail,
        closeOutreach,
    } = useBulkActions({ selectedIds, clearSelection });

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
                                    checked={isAllSelected}
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
                            <CandidateRow
                                key={candidate.id}
                                candidate={candidate}
                                isSelected={selectedIds.has(candidate.id)}
                                onToggleSelection={toggleSelection}
                            />
                        ))}
                    </TableBody>
                </Table>
            </div>

            <BulkActionBar
                selectedCount={selectedIds.size}
                onClear={clearSelection}
                onRejectAll={handleBulkReject}
                onMoveTo={handleBulkMove}
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
