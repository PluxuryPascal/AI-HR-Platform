"use client";

import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Briefcase, Users, Clock } from "lucide-react";
import { useTranslations } from "next-intl";

import { useGetDashboardStats } from "@/features/dashboard/api/use-dashboard";
import { Loader2 } from "lucide-react";

export function JobsStats() {
    const t = useTranslations("Stats");
    const { data: stats, isLoading } = useGetDashboardStats();

    if (isLoading) {
        return (
            <div className="grid gap-4 md:grid-cols-3">
                {[1, 2, 3].map((i) => (
                    <Card key={i} className="animate-pulse">
                        <CardHeader className="pb-2">
                            <div className="h-4 w-24 bg-muted rounded" />
                        </CardHeader>
                        <CardContent>
                            <div className="h-8 w-12 bg-muted rounded mb-2" />
                            <div className="h-3 w-32 bg-muted rounded" />
                        </CardContent>
                    </Card>
                ))}
            </div>
        );
    }

    return (
        <div className="grid gap-4 md:grid-cols-3">
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                        {t("activeJobs")}
                    </CardTitle>
                    <Briefcase className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{stats?.active_jobs || 0}</div>
                    <p className="text-xs text-muted-foreground">
                        {t("fromLastMonth", { count: stats?.active_jobs_delta || 0 })}
                    </p>
                </CardContent>
            </Card>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                        {t("totalCandidates")}
                    </CardTitle>
                    <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{stats?.total_candidates || 0}</div>
                    <p className="text-xs text-muted-foreground">
                        {t("fromLastMonth", { count: stats?.total_candidates_delta || 0 })}
                    </p>
                </CardContent>
            </Card>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                        {t("pendingScreenings")}
                    </CardTitle>
                    <Clock className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold text-amber-500">{stats?.upcoming_interviews || 0}</div>
                    <p className="text-xs text-muted-foreground">
                        {t("requiresAttention")}
                    </p>
                </CardContent>
            </Card>
        </div>
    );
}
