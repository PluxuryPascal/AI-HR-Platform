import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Plus } from "lucide-react";
import { JobsStats } from "@/features/jobs/components/jobs-stats";
import { JobsTable } from "@/features/jobs/components/jobs-table";
import { useTranslations } from "next-intl";

export default function JobsPage() {
    const t = useTranslations("Jobs");

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">{t("title")}</h2>
                <div className="flex items-center space-x-2">
                    <Button>
                        <Plus className="mr-2 h-4 w-4" /> {t("create")}
                    </Button>
                </div>
            </div>
            <div className="space-y-4">
                <JobsStats />
                <div className="flex items-center justify-between">
                    <div className="flex flex-1 items-center space-x-2">
                        <Input
                            placeholder={t("filterPlaceholder")}
                            className="h-8 w-[150px] lg:w-[250px]"
                        />
                        <Select>
                            <SelectTrigger className="h-8 w-[180px]">
                                <SelectValue placeholder={t("statusPlaceholder")} />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">{t("statusAll")}</SelectItem>
                                <SelectItem value="active">{t("statusActive")}</SelectItem>
                                <SelectItem value="closed">{t("statusClosed")}</SelectItem>
                                <SelectItem value="draft">{t("statusDraft")}</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                </div>
                <JobsTable />
            </div>
        </div>
    );
}
