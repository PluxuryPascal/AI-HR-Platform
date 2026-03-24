"use client";
import { useState } from "react";

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal, FileText, Archive, Eye, Briefcase, Users } from "lucide-react";
import { useTranslations } from "next-intl";
import { EmptyState } from "@/components/shared/empty-state";
import { useRouter } from "@/i18n/routing";
import { Pagination } from "@/components/shared/pagination";
import { ChevronUp, ChevronDown, ChevronsUpDown } from "lucide-react";


import { useGetJobs } from "../api/use-get-jobs";
import { usePublishJob, useCloseJob, useArchiveJob } from "../api/use-job-actions";
import { JobStatus, JobFilter } from "../types/job";
import { Loader2, Zap, XCircle, ArchiveRestore } from "lucide-react";

interface JobsTableProps {
    filter?: JobFilter;
}

interface StatusBadgeProps {
    status: JobStatus;
}

const StatusBadge = ({ status }: StatusBadgeProps) => {
    const t = useTranslations("Jobs");

    switch (status) {
        case "status_published":
            return (
                <Badge className="bg-green-100 text-green-700 hover:bg-green-100/80 border-green-200 shadow-none dark:bg-green-900/30 dark:text-green-400 dark:border-green-800">
                    {t("statusActive")}
                </Badge>
            );
        case "status_closed":
            return (
                <Badge variant="secondary" className="bg-slate-100 text-slate-700 hover:bg-slate-100/80 border-slate-200 shadow-none dark:bg-slate-800 dark:text-slate-400 dark:border-slate-800">
                    {t("statusClosed")}
                </Badge>
            );
        case "status_draft":
            return (
                <Badge variant="outline" className="text-slate-500 border-slate-300 dark:text-slate-400 dark:border-slate-700">
                    {t("statusDraft")}
                </Badge>
            );
        case "status_archived":
            return (
                <Badge variant="outline" className="text-slate-400 border-slate-200">
                    {t("statusArchived") || "Archived"}
                </Badge>
            );
        default:
            return <Badge variant="outline">{status}</Badge>;
    }
};

export function JobsTable({ filter }: JobsTableProps) {
    const t = useTranslations("Jobs");
    const router = useRouter();
    
    const [page, setPage] = useState(1);
    const [perPage] = useState(10);
    const [sortId, setSortId] = useState<string>("created_at");
    const [sortDesc, setSortDesc] = useState(true);

    const { data: jobsResponse, isLoading, error } = useGetJobs({ 
        page, 
        per_page: perPage, 
        filter: {
            ...filter,
            sort: {
                sort_id: sortId,
                sort_desc: sortDesc
            }
        } 
    });
    const publishJob = usePublishJob();
    const closeJob = useCloseJob();
    const archiveJob = useArchiveJob();

    const handleSort = (id: string) => {
        if (sortId === id) {
            setSortDesc(!sortDesc);
        } else {
            setSortId(id);
            setSortDesc(true);
        }
        setPage(1); // Reset to first page on sort
    };

    const SortIcon = ({ id }: { id: string }) => {
        if (sortId !== id) return <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />;
        return sortDesc ? <ChevronDown className="ml-2 h-4 w-4" /> : <ChevronUp className="ml-2 h-4 w-4" />;
    };

    if (isLoading) {
        return (
            <div className="flex justify-center p-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-8 text-center text-red-500">
                {t("error") || "Failed to load jobs"}
            </div>
        );
    }

    const jobs = jobsResponse?.data || [];

    if (jobs.length === 0) {
        return (
            <EmptyState
                icon={Briefcase}
                title={t("empty.title")}
                description={t("empty.desc")}
                actionLabel={t("empty.button")}
                onAction={() => router.push("/dashboard/jobs/new")}
            />
        );
    }

    return (
        <div className="rounded-md border bg-card">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[30%] cursor-pointer hover:bg-muted/50 transition-colors" onClick={() => handleSort("title")}>
                            <div className="flex items-center">
                                {t("table.title")}
                                <SortIcon id="title" />
                            </div>
                        </TableHead>
                        <TableHead>{t("table.department")}</TableHead>
                        <TableHead>{t("table.type")}</TableHead>
                        <TableHead>{t("table.salary")}</TableHead>
                        <TableHead className="cursor-pointer hover:bg-muted/50 transition-colors" onClick={() => handleSort("created_at")}>
                            <div className="flex items-center">
                                {t("table.created")}
                                <SortIcon id="created_at" />
                            </div>
                        </TableHead>
                        <TableHead>{t("table.status")}</TableHead>
                        <TableHead className="text-right">{t("table.actions")}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {jobs.map((job) => (
                        <TableRow key={job.id} className="hover:bg-muted/50 transition-colors">
                            <TableCell className="font-medium">{job.title}</TableCell>
                            <TableCell>{job.department_name || "-"}</TableCell>
                            <TableCell className="capitalize">{job.work_format}</TableCell>
                            <TableCell>
                                {job.salary_min || job.salary_max ? (
                                    <>
                                        {job.salary_min && `${job.salary_min.toLocaleString()}`}
                                        {job.salary_min && job.salary_max && " - "}
                                        {job.salary_max && `${job.salary_max.toLocaleString()}`}
                                        {` ${job.currency}`}
                                    </>
                                ) : (
                                    "-"
                                )}
                            </TableCell>
                            <TableCell>{new Date(job.created_at).toLocaleDateString()}</TableCell>
                            <TableCell>
                                <StatusBadge status={job.status} />
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
                                        <DropdownMenuLabel>{t("table.actions")}</DropdownMenuLabel>
                                        <DropdownMenuItem onClick={() => router.push(`/dashboard/candidates?job_id=${job.id}`)}>
                                            <Users className="mr-2 h-4 w-4" />
                                            {t("table.viewCandidates")}
                                        </DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => router.push(`/dashboard/jobs/${job.id}/edit`)}>
                                            <FileText className="mr-2 h-4 w-4" />
                                            {t("table.editJob")}
                                        </DropdownMenuItem>
                                        <DropdownMenuSeparator />
                                        
                                        {job.status !== "status_published" && (
                                            <DropdownMenuItem onClick={() => publishJob.mutate(job.id)} disabled={publishJob.isPending}>
                                                <Zap className="mr-2 h-4 w-4 text-amber-500" />
                                                Опубликовать
                                            </DropdownMenuItem>
                                        )}
                                        
                                        {job.status !== "status_closed" && job.status !== "status_draft" && (
                                            <DropdownMenuItem onClick={() => closeJob.mutate(job.id)} disabled={closeJob.isPending}>
                                                <XCircle className="mr-2 h-4 w-4 text-red-500" />
                                                Закрыть
                                            </DropdownMenuItem>
                                        )}
                                        
                                        {job.status !== "status_archived" && job.status !== "status_draft" && (
                                            <DropdownMenuItem onClick={() => archiveJob.mutate(job.id)} disabled={archiveJob.isPending}>
                                                <Archive className="mr-2 h-4 w-4 text-slate-500" />
                                                Архивировать
                                            </DropdownMenuItem>
                                        )}
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
            <Pagination 
                currentPage={page} 
                totalPages={jobsResponse?.meta?.total_pages || 1} 
                onPageChange={setPage} 
            />
        </div>
    );
}
