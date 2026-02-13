import { useTranslations } from "next-intl";
import { JobForm } from "@/features/jobs/components/job-form";
import { Button } from "@/components/ui/button";
import { ChevronLeft } from "lucide-react";
import { Link } from "@/i18n/routing";


export default function CreateJobPage() {
    const t = useTranslations("JobWizard");

    return (
        <div className="relative min-h-screen">
            <div className="container relative z-10 max-w-3xl py-10">
                <div className="mb-6">
                    <Link href="/dashboard/jobs">
                        <Button variant="ghost" className="pl-0 hover:pl-2 transition-all gap-2 text-muted-foreground hover:text-foreground">
                            <ChevronLeft className="h-4 w-4" />
                            {t("backToJobs")}
                        </Button>
                    </Link>
                    <div className="mt-4">
                        <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
                        <p className="text-muted-foreground mt-2">
                            {t("subtitle")}
                        </p>
                    </div>
                </div>

                <JobForm />
            </div>
        </div>
    );
}
