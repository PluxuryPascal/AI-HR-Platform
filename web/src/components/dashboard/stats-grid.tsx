"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Users, Briefcase, Clock, CalendarCheck } from "lucide-react";
import { glassCard } from "@/lib/styles";
import { useGetDashboardStats } from "@/features/dashboard/api/use-dashboard";

export function StatsGrid() {
    const t = useTranslations("Dashboard.stats");
    const tStats = useTranslations("Stats");

    const { data: dashboardStats, isLoading } = useGetDashboardStats();

    const stats = [
        {
            title: t("applicants"),
            value: dashboardStats?.total_candidates?.toLocaleString() || "0",
            icon: Users,
            trend: dashboardStats?.total_candidates_delta ? `${dashboardStats.total_candidates_delta > 0 ? "+" : ""}${dashboardStats.total_candidates_delta}%` : "0%",
            trendUp: (dashboardStats?.total_candidates_delta || 0) >= 0,
        },
        {
            title: t("activeJobs"),
            value: dashboardStats?.active_jobs?.toString() || "0",
            icon: Briefcase,
            trend: dashboardStats?.active_jobs_delta ? `${dashboardStats.active_jobs_delta > 0 ? "+" : ""}${dashboardStats.active_jobs_delta}` : "0",
            trendUp: (dashboardStats?.active_jobs_delta || 0) >= 0,
        },
        {
            title: t("interviews"),
            value: dashboardStats?.upcoming_interviews?.toString() || "0",
            icon: CalendarCheck,
            trend: dashboardStats?.interviews_delta ? `${dashboardStats.interviews_delta > 0 ? "+" : ""}${dashboardStats.interviews_delta}` : "0",
            trendUp: (dashboardStats?.interviews_delta || 0) >= 0,
        },
        {
            title: t("timeToHire"),
            value: `${dashboardStats?.avg_time_to_hire_days?.toFixed(1) || "0"} ${tStats("days")}`,
            icon: Clock,
            trend: dashboardStats?.avg_time_to_hire_delta ? `${dashboardStats.avg_time_to_hire_delta > 0 ? "+" : ""}${dashboardStats.avg_time_to_hire_delta.toFixed(1)} ${tStats("days")}` : `0 ${tStats("days")}`,
            trendUp: (dashboardStats?.avg_time_to_hire_delta || 0) <= 0, // Lower is better for time to hire
        },
    ];

    if (isLoading) {
        return (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {[1, 2, 3, 4].map((i) => (
                    <Card key={i} className={`${glassCard} animate-pulse h-32`} />
                ))}
            </div>
        );
    }

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {stats.map((stat, index) => (
                <Card key={index} className={`${glassCard} hover:-translate-y-1 group`}>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">
                            {stat.title}
                        </CardTitle>
                        <div className="p-2 rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
                            <stat.icon className="h-4 w-4" />
                        </div>
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stat.value}</div>
                        <p className="text-xs text-muted-foreground">
                            <span className={stat.trendUp ? "text-emerald-500 font-medium" : "text-rose-500 font-medium"}>
                                {stat.trend}
                            </span>{" "}
                            {tStats("fromLastMonth", { count: "" }).replace("+ ", "")}
                        </p>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}
