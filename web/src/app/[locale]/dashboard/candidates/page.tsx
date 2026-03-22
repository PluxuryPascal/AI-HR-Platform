"use client";
import { useTranslations } from "next-intl";
import { Plus } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

import { CandidatesTable } from "@/features/candidates/components/candidates-table";
import { UploadResumeDialog } from "@/features/candidates/components/upload-resume-dialog";

import { useState, useEffect } from "react";
import { useGetJobs } from "@/features/jobs/api/use-get-jobs";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

export default function CandidatesPage() {
    const t = useTranslations("Candidates");
    const { data: jobsResponse } = useGetJobs({ per_page: 100 });
    const jobs = jobsResponse?.data || [];
    
    const [selectedJobId, setSelectedJobId] = useState<string>("");

    // Set first job as default if none selected
    useEffect(() => {
        if (!selectedJobId && jobs.length > 0) {
            setSelectedJobId(jobs[0].id);
        }
    }, [jobs, selectedJobId]);

    return (
        <div className="flex flex-col gap-6 p-6">
            <div className="flex items-center justify-between">
                <div className="space-y-1">
                    <h2 className="text-2xl font-bold tracking-tight">{t("title")}</h2>
                    <p className="text-muted-foreground">
                        Manage and track your candidate pipeline by job.
                    </p>
                </div>
                <div className="flex items-center gap-4">
                    <div className="w-[250px]">
                        <Select value={selectedJobId} onValueChange={setSelectedJobId}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select a job" />
                            </SelectTrigger>
                            <SelectContent>
                                {jobs.map((job) => (
                                    <SelectItem key={job.id} value={job.id}>
                                        {job.title}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                    <UploadResumeDialog jobId={selectedJobId}>
                        <Button disabled={!selectedJobId}>
                            <Plus className="mr-2 h-4 w-4" />
                            {t("uploadBtn")}
                        </Button>
                    </UploadResumeDialog>
                </div>
            </div>
            <Separator />
            <CandidatesTable jobId={selectedJobId} />
        </div>
    );
}
