"use client";

import { useTranslations } from "next-intl";
import { OverviewChart } from "@/components/dashboard/overview-chart";
import { RecentActivity } from "@/components/dashboard/recent-activity";
import { StatsGrid } from "@/components/dashboard/stats-grid";

export default function DashboardPage() {
    const t = useTranslations("Dashboard");

    return (
        <div className="flex-1 space-y-8 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">
                    {t("welcome", { name: "Recruiter" })}
                </h2>
            </div>
            <div className="space-y-4">
                <StatsGrid />
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                    <div className="col-span-1 md:col-span-2 lg:col-span-4">
                        <OverviewChart />
                    </div>
                    <div className="col-span-1 md:col-span-2 lg:col-span-3">
                        <RecentActivity />
                    </div>
                </div>
            </div>
        </div>
    );
}
