"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Users, Briefcase, Clock, CalendarCheck } from "lucide-react";
import { glassCard } from "@/lib/styles";

export function StatsGrid() {
    const t = useTranslations("Dashboard.stats");
    const tStats = useTranslations("Stats");

    const stats = [
        {
            title: t("applicants"),
            value: "1,240",
            icon: Users,
            trend: "+12%",
            trendUp: true,
        },
        {
            title: t("activeJobs"),
            value: "8",
            icon: Briefcase,
            trend: "+2",
            trendUp: true,
        },
        {
            title: t("interviews"),
            value: "12",
            icon: CalendarCheck,
            trend: "+4",
            trendUp: true,
        },
        {
            title: t("timeToHire"),
            value: "12 Days",
            icon: Clock,
            trend: "-3 Days",
            trendUp: true, // actually good in this context
        },
    ];

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
