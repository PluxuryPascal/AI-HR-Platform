"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { KanbanBoard } from "@/features/screening/components/kanban-board";
import { useJobsStore } from "@/store/use-jobs-store";
import { 
    Select as UISelect, 
    SelectContent as UISelectContent, 
    SelectItem as UISelectItem, 
    SelectTrigger as UISelectTrigger, 
    SelectValue as UISelectValue 
} from "@/components/ui/select";

export default function ScreeningPage() {
    const t = useTranslations("Screening");
    const jobs = useJobsStore((state) => state.jobs);
    const [selectedJobId, setSelectedJobId] = useState<string>(jobs[0]?.id || "");

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
