"use client";

import { useTranslations } from "next-intl";
import { OverviewChart } from "@/components/dashboard/overview-chart";
import { RecentActivity } from "@/components/dashboard/recent-activity";
import { StatsGrid } from "@/components/dashboard/stats-grid";

export default function DashboardPage() {
    const t = useTranslations("Dashboard");

    return (
        <div className="relative flex-1 min-h-full overflow-hidden">
            {/* Background Aura */}
            <div className="absolute inset-0 z-0 pointer-events-none overflow-hidden">
                <div className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full bg-indigo-500/20 blur-[120px] mix-blend-screen animate-pulse dark:bg-indigo-500/10" />
                <div className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full bg-violet-500/20 blur-[120px] mix-blend-screen animate-pulse delay-1000 dark:bg-violet-500/10" />
            </div>

            <div className="relative z-10 space-y-6 p-6 pt-6">
                <div className="flex items-center justify-between space-y-2">
                    <h2 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-foreground to-foreground/60">
                        {t("welcome", { name: "Recruiter" })}
                    </h2>
                </div>
                <div className="space-y-6">
                    <StatsGrid />
                    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-7">
                        <div className="col-span-1 md:col-span-2 lg:col-span-4">
                            <OverviewChart />
                        </div>
                        <div className="col-span-1 md:col-span-2 lg:col-span-3">
                            <RecentActivity />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
