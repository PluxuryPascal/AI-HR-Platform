"use client";

import { useTranslations } from "next-intl";
import {
    Area,
    Bar,
    ComposedChart,
    CartesianGrid,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { WidgetErrorBoundary } from "@/components/shared/widget-error-boundary";
import { useGetDashboardDynamics } from "@/features/dashboard/api/use-dashboard";
import { useMemo } from "react";


import { glassCard } from "@/lib/styles";

export function OverviewChart() {
    const t = useTranslations("Dashboard.charts");

    const endDate = useMemo(() => new Date().toISOString().split('T')[0], []);
    const startDate = useMemo(() => {
        const d = new Date();
        d.setMonth(d.getMonth() - 6);
        return d.toISOString().split('T')[0];
    }, []);

    const { data: rawData = [], isLoading } = useGetDashboardDynamics(startDate, endDate, "monthly");

    const chartData = useMemo(() => {
        const monthFormatter = new Intl.DateTimeFormat('ru-RU', { month: 'short' });
        return (rawData || []).map(d => ({
            name: monthFormatter.format(new Date(d.date)).replace('.', ''),
            total: d.count
        }));
    }, [rawData]);

    return (
        <Card className={`col-span-1 md:col-span-2 ${glassCard}`}>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    {t("applicationsOverTime")}
                </CardTitle>
            </CardHeader>
            <CardContent className="pl-2">
                <WidgetErrorBoundary>
                    <ResponsiveContainer width="100%" height={350}>
                        <ComposedChart data={isLoading ? [] : chartData}>
                            <defs>
                                <linearGradient id="barGradient" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="var(--chart-gradient-start)" />
                                    <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity="var(--chart-gradient-end)" />
                                </linearGradient>
                                <linearGradient id="areaGradient" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity={0.2} />
                                    <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0} />
                                </linearGradient>
                            </defs>
                            <CartesianGrid
                                strokeDasharray="3 3"
                                vertical={false}
                                stroke="hsl(var(--border))"
                                opacity={0.4}
                            />
                            <XAxis
                                dataKey="name"
                                stroke="hsl(var(--muted-foreground))"
                                fontSize={12}
                                tickLine={false}
                                axisLine={false}
                                dy={10}
                            />
                            <YAxis
                                stroke="hsl(var(--muted-foreground))"
                                fontSize={12}
                                tickLine={false}
                                axisLine={false}
                                tickFormatter={(value) => `${value}`}
                                dx={-10}
                            />
                            <Tooltip
                                cursor={{ fill: "hsl(var(--muted) / 0.2)", radius: 4 }}
                                content={({ active, payload, label }) => {
                                    if (active && payload && payload.length) {
                                        return (
                                            <div className="overflow-hidden rounded-xl border border-border/50 bg-white dark:bg-slate-950 p-0 shadow-xl">
                                                <div className="border-b border-white/5 bg-white/5 p-3">
                                                    <span className="text-sm font-medium text-foreground">
                                                        {label}
                                                    </span>
                                                </div>
                                                <div className="p-3">
                                                    <div className="flex items-center justify-between gap-8">
                                                        <div className="flex items-center gap-2">
                                                            <div className="h-2 w-2 rounded-full bg-primary ring-2 ring-primary/20" />
                                                            <span className="text-xs text-muted-foreground">Applications</span>
                                                        </div>
                                                        <span className="font-bold tabular-nums">
                                                            {payload[0].value}
                                                        </span>
                                                    </div>
                                                </div>
                                            </div>
                                        );
                                    }
                                    return null;
                                }}
                            />



                            <Bar
                                dataKey="total"
                                fill="url(#barGradient)"
                                stroke="currentColor"
                                strokeOpacity={0.2}
                                radius={[6, 6, 0, 0]}
                                isAnimationActive={true}
                                animationDuration={1500}
                                animationEasing="ease-in-out"
                                className="text-primary"
                            />
                        </ComposedChart>
                    </ResponsiveContainer>
                </WidgetErrorBoundary>
            </CardContent>
        </Card>
    );
}
