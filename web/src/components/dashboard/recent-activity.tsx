"use client";

import { useTranslations } from "next-intl";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Briefcase, UserCheck, MessageSquare, Clock, Sparkles } from "lucide-react";
import { glassCard } from "@/lib/styles";
import { useGetRecentActivity } from "@/features/dashboard/api/use-dashboard";

export function RecentActivity() {
    const t = useTranslations("Dashboard.activity");
    const { data: rawActivities = [], isLoading } = useGetRecentActivity(10);

    const formatRelativeTime = (date: Date) => {
        const now = new Date();
        const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

        if (diffInSeconds < 60) return "только что";
        const diffInMinutes = Math.floor(diffInSeconds / 60);
        if (diffInMinutes < 60) return `${diffInMinutes} мин. назад`;
        const diffInHours = Math.floor(diffInMinutes / 60);
        if (diffInHours < 24) return `${diffInHours} ч. назад`;
        const diffInDays = Math.floor(diffInHours / 24);
        return `${diffInDays} дн. назад`;
    };

    const activities = (rawActivities || []).map(log => {
        let icon = Clock;
        let color = "text-muted-foreground";
        let action = log.action_code;
        let target = log.job_title || log.candidate_first_name || "Object";
        let isAi = log.actor_type === "ai";

        if (log.action_code.includes("screening")) {
            icon = Sparkles;
            color = "text-indigo-500";
        } else if (log.action_code.includes("job")) {
            icon = Briefcase;
            color = "text-blue-500";
        } else if (log.action_code.includes("interview")) {
            icon = MessageSquare;
            color = "text-purple-500";
        }

        return {
            id: log.id,
            type: log.action_code,
            user: log.actor_name || "System",
            action: action,
            target: target,
            meta: log.match_score ? `${log.match_score}% Match` : "",
            time: formatRelativeTime(new Date(log.created_at)),
            icon: icon,
            color: color,
            isAi: isAi,
        };
    });

    return (
        <Card className={`col-span-1 ${glassCard}`}>
            <CardHeader>
                <CardTitle>{t("title")}</CardTitle>
            </CardHeader>
            <CardContent>
                <ScrollArea className="h-[350px] pr-4">
                    <div className="space-y-4">
                        {activities.length === 0 ? (
                            <p className="text-sm text-muted-foreground text-center py-4">{t("empty")}</p>
                        ) : (
                            activities.map((activity) => (
                                <div
                                    key={activity.id}
                                    className={`flex items-start gap-4 p-3 rounded-xl transition-all ${activity.isAi
                                        ? "bg-indigo-500/5 border border-indigo-500/10 hover:bg-indigo-500/10"
                                        : "hover:bg-white/5 border border-transparent"
                                        }`}
                                >
                                    <div className={`mt-0.5 rounded-full p-2 ${activity.isAi ? "bg-indigo-500/10 text-indigo-500" : "bg-muted text-muted-foreground"}`}>
                                        <activity.icon className="h-4 w-4" />
                                    </div>
                                    <div className="space-y-1 flex-1">
                                        <p className="text-sm font-medium leading-none flex flex-wrap gap-1 items-center">
                                            <span className="font-semibold">{activity.user}</span>
                                            <span className="text-muted-foreground">{activity.action}</span>
                                            <span className="font-semibold">{activity.target}</span>
                                        </p>
                                        <div className="flex items-center gap-2">
                                            <div className="flex items-center text-xs text-muted-foreground">
                                                <Clock className="mr-1 h-3 w-3" />
                                                {activity.time}
                                            </div>
                                            {activity.meta && (
                                                <span className={`text-[10px] uppercase font-bold px-1.5 py-0.5 rounded ${activity.isAi ? "bg-green-500/10 text-green-500 animate-pulse" : "bg-muted text-muted-foreground"
                                                    }`}>
                                                    {activity.meta}
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            )))}
                    </div>
                </ScrollArea>
            </CardContent>
        </Card>
    );
}
