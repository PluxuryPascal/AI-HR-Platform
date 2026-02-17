import { useTranslations } from "next-intl";
import { Plus } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

import { CandidatesTable } from "@/features/candidates/components/candidates-table";
import { UploadResumeDialog } from "@/features/candidates/components/upload-resume-dialog";

export default function CandidatesPage() {
    const t = useTranslations("Candidates");

    return (
        <div className="flex flex-col gap-6 p-6">
            <div className="flex items-center justify-between">
                <div className="space-y-1">
                    <h2 className="text-2xl font-bold tracking-tight">{t("title")}</h2>
                    <p className="text-muted-foreground">
                        Manage and track your candidate pipeline.
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    <UploadResumeDialog>
                        <Button>
                            <Plus className="mr-2 h-4 w-4" />
                            {t("uploadBtn")}
                        </Button>
                    </UploadResumeDialog>
                </div>
            </div>
            <Separator />
            <CandidatesTable />
        </div>
    );
}
