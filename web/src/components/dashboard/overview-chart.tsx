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

const data = [
    { name: "Jan", total: 124, predicted: 140 },
    { name: "Feb", total: 145, predicted: 160 },
    { name: "Mar", total: 112, predicted: 130 },
    { name: "Apr", total: 178, predicted: 190 },
    { name: "May", total: 234, predicted: 220 },
    { name: "Jun", total: 215, predicted: 240 },
];

import { glassCard } from "@/lib/styles";

export function OverviewChart() {
    const t = useTranslations("Dashboard.charts");

    return (
        <Card className={`col-span-1 md:col-span-2 ${glassCard}`}>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    {t("applicationsOverTime")}
                    <span className="inline-flex items-center rounded-full border border-sky-500/30 bg-sky-500/10 px-2 py-0.5 text-xs font-medium text-sky-500">
                        AI Forecast
                    </span>
                </CardTitle>
            </CardHeader>
            <CardContent className="pl-2">
                <WidgetErrorBoundary>
                    <ResponsiveContainer width="100%" height={350}>
                        <ComposedChart data={data}>
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
                                                <div className="p-3 space-y-3">
                                                    <div className="flex items-center justify-between gap-8">
                                                        <div className="flex items-center gap-2">
                                                            <div className="h-2 w-2 rounded-full bg-primary ring-2 ring-primary/20" />
                                                            <span className="text-xs text-muted-foreground">Applications</span>
                                                        </div>
                                                        <span className="font-bold tabular-nums">
                                                            {payload[0].value}
                                                        </span>
                                                    </div>
                                                    <div className="flex items-center justify-between gap-8">
                                                        <div className="flex items-center gap-2">
                                                            <div className="h-2 w-2 rounded-full bg-sky-500 ring-2 ring-sky-500/20" />
                                                            <span className="text-xs text-muted-foreground">AI Forecast</span>
                                                        </div>
                                                        <span className="font-bold tabular-nums text-sky-500">
                                                            {payload[1]?.value}
                                                        </span>
                                                    </div>
                                                </div>
                                                <div className="bg-gradient-to-r from-emerald-500/10 via-teal-500/10 to-emerald-500/10 p-2 px-3">
                                                    <div className="flex items-center gap-2">
                                                        <span className="relative flex h-2 w-2">
                                                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                                                            <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                                                        </span>
                                                        <span className="text-[10px] font-medium text-emerald-500 uppercase tracking-wider">
                                                            Projected Growth: +12%
                                                        </span>
                                                    </div>
                                                </div>
                                            </div>
                                        );
                                    }
                                    return null;
                                }}
                            />

                            <Area
                                type="monotone"
                                dataKey="predicted"
                                stroke="hsl(var(--primary))"
                                strokeWidth={2}
                                fill="url(#areaGradient)"
                                className="opacity-50"
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
