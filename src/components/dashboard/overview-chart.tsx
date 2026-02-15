"use client";

import { useTranslations } from "next-intl";
import {
    Bar,
    BarChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { WidgetErrorBoundary } from "@/components/shared/widget-error-boundary";

const data = [
    { name: "Jan", total: 124 },
    { name: "Feb", total: 145 },
    { name: "Mar", total: 112 },
    { name: "Apr", total: 178 },
    { name: "May", total: 234 },
    { name: "Jun", total: 215 },
];

export function OverviewChart() {
    const t = useTranslations("Dashboard.charts");

    return (
        <Card className="col-span-1 md:col-span-2">
            <CardHeader>
                <CardTitle>{t("applicationsOverTime")}</CardTitle>
            </CardHeader>
            <CardContent className="pl-2">
                <WidgetErrorBoundary>
                    <ResponsiveContainer width="100%" height={350}>
                        <BarChart data={data}>
                            <XAxis
                                dataKey="name"
                                stroke="#888888"
                                fontSize={12}
                                tickLine={false}
                                axisLine={false}
                            />
                            <YAxis
                                stroke="#888888"
                                fontSize={12}
                                tickLine={false}
                                axisLine={false}
                                tickFormatter={(value) => `${value}`}
                            />
                            <Tooltip
                                cursor={{ fill: "transparent" }}
                                content={({ active, payload, label }) => {
                                    if (active && payload && payload.length) {
                                        return (
                                            <div className="rounded-lg border bg-background p-2 shadow-sm">
                                                <div className="grid grid-cols-2 gap-2">
                                                    <div className="flex flex-col">
                                                        <span className="text-[0.70rem] uppercase text-muted-foreground">
                                                            {label}
                                                        </span>
                                                        <span className="font-bold text-muted-foreground">
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
                                fill="currentColor"
                                radius={[4, 4, 0, 0]}
                                className="fill-primary"
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </WidgetErrorBoundary>
            </CardContent>
        </Card>
    );
}
