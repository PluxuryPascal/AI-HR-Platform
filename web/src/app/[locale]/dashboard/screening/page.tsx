"use client";

import { useState, useEffect } from "react";
import { useTranslations } from "next-intl";
import { useGetJobs } from "@/features/jobs/api/use-get-jobs";
import { KanbanBoard } from "@/features/screening/components/kanban-board";
import { 
    Select as UISelect, 
    SelectContent as UISelectContent, 
    SelectItem as UISelectItem, 
    SelectTrigger as UISelectTrigger, 
    SelectValue as UISelectValue 
} from "@/components/ui/select";

export default function ScreeningPage() {
    const t = useTranslations("Screening");
    const { data: jobsResponse } = useGetJobs();
    const jobs = jobsResponse?.data || [];
    const [selectedJobId, setSelectedJobId] = useState<string>("");

    // Set initial selected job once jobs are loaded
    useEffect(() => {
        if (!selectedJobId && jobs.length > 0) {
            setSelectedJobId(jobs[0].id);
        }
    }, [jobs, selectedJobId]);

    return (
        <div className="flex flex-col h-[calc(100vh-4rem)]">
            <div className="flex items-center justify-between px-6 py-4 border-b">
                <div className="flex items-center gap-4">
                    <h1 className="text-2xl font-bold tracking-tight">{t("title")}</h1>
                    <UISelect value={selectedJobId} onValueChange={setSelectedJobId}>
                        <UISelectTrigger className="w-[280px]">
                            <UISelectValue placeholder={t("jobSelectorPlaceholder")} />
                        </UISelectTrigger>
                        <UISelectContent>
                            {jobs.map((job) => (
                                <UISelectItem key={job.id} value={job.id}>
                                    {job.title}
                                </UISelectItem>
                            ))}
                        </UISelectContent>
                    </UISelect>
                </div>
            </div>
            <div className="flex-1 overflow-hidden p-6">
                {selectedJobId ? (
                    <KanbanBoard jobId={selectedJobId} />
                ) : (
                    <div className="flex items-center justify-center h-full text-muted-foreground">
                        {t("noJobsSelected")}
                    </div>
                )}
            </div>
        </div>
    );
}
