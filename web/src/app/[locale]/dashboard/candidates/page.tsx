"use client";
import { useTranslations } from "next-intl";
import { Plus } from "lucide-react";

import { useGetJobs } from "@/features/jobs/api/use-get-jobs";
import { useGetJob } from "@/features/jobs/api/use-get-job";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";

import { CandidatesTable } from "@/features/candidates/components/candidates-table";
import { UploadResumeDialog } from "@/features/candidates/components/upload-resume-dialog";

import { useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

export default function CandidatesPage() {
    const t = useTranslations("Candidates");
    const searchParams = useSearchParams();
    const jobIdFromUrl = searchParams.get("job_id");
    const [isJobsEnabled, setIsJobsEnabled] = useState(false);

    const { data: jobsResponse } = useGetJobs({ 
        per_page: 100, 
        enabled: !jobIdFromUrl || isJobsEnabled,
        refetchInterval: false 
    });
    const { data: jobResponse } = useGetJob(jobIdFromUrl || "");
    
    const jobs = jobsResponse?.data || [];
    const currentJob = jobResponse;
    
    const [selectedJobId, setSelectedJobId] = useState<string>(jobIdFromUrl || "");
    const [search, setSearch] = useState("");
    const [debouncedSearch, setDebouncedSearch] = useState("");

    // Debounce search
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(search);
        }, 500);
        return () => clearTimeout(timer);
    }, [search]);

    // Set default job if none selected and jobs loaded
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
                    <Input
                        placeholder="Search by name or email..."
                        className="h-9 w-[250px]"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                    />
                    <div className="w-[250px]">
                        <Select 
                            value={selectedJobId} 
                            onValueChange={setSelectedJobId}
                            onOpenChange={(open) => open && setIsJobsEnabled(true)}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder={currentJob?.title || "Select a job"} />
                            </SelectTrigger>
                            <SelectContent>
                                {jobs.length > 0 ? (
                                    jobs.map((job) => (
                                        <SelectItem key={job.id} value={job.id}>
                                            {job.title}
                                        </SelectItem>
                                    ))
                                ) : currentJob ? (
                                    <SelectItem value={currentJob.id}>
                                        {currentJob.title}
                                    </SelectItem>
                                ) : null}
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
            {selectedJobId && (
                <CandidatesTable 
                    jobId={selectedJobId} 
                    filter={{ first_name: debouncedSearch || undefined }} 
                />
            )}
        </div>
    );
}
