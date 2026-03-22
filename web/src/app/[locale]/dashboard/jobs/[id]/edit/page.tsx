"use client";

import { useTranslations } from "next-intl";
import { Link } from "@/i18n/routing";
import { JobForm } from "@/features/jobs/components/job-form";
import { useGetJob } from "@/features/jobs/api/use-get-job";
import { Loader2, ChevronLeft } from "lucide-react";
import { use } from "react";

export default function EditJobPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = use(params);
    const t = useTranslations("JobWizard");
    const { data: job, isLoading, error } = useGetJob(id);

    if (isLoading) {
        return (
            <div className="flex items-center justify-center p-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        );
    }

    if (error || !job) {
        return (
            <div className="p-8 text-center text-red-500">
                Failed to load job data.
            </div>
        );
    }

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between">
                <div className="space-y-1">
                    <div className="flex items-center gap-2">
                        <Link href="/dashboard/jobs" className="text-muted-foreground hover:text-primary transition-colors">
                            <ChevronLeft className="h-4 w-4" />
                        </Link>
                        <h2 className="text-3xl font-bold tracking-tight">
                            {t("titleEdit") || "Редактирование вакансии"}
                        </h2>
                    </div>
                    <p className="text-muted-foreground">
                        {job.title}
                    </p>
                </div>
            </div>

            <JobForm id={id} initialData={job} />
        </div>
    );
}
