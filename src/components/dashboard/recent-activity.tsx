"use client";

import { useTranslations } from "next-intl";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Briefcase, UserCheck, MessageSquare, Clock, Sparkles } from "lucide-react";

export function RecentActivity() {
    const t = useTranslations("Dashboard.activity");

    const activities = [
        {
            type: "screening",
            user: "AI Agent",
            action: "finished screening",
            target: "Alex Doe",
            meta: "92% Match",
            time: "10m ago",
            icon: Sparkles,
            color: "text-indigo-500",
            isAi: true,
        },
        {
            type: "job",
            user: "System",
            action: "New Job",
            target: "Senior Backend",
            meta: "created",
            time: "2h ago",
            icon: Briefcase,
            color: "text-blue-500",
        },
        {
            type: "interview",
            user: "Maria Garcia",
            action: "moved to",
            target: "Interview",
            meta: "",
            time: "5h ago",
            icon: MessageSquare,
            color: "text-purple-500",
        },
        {
            type: "screening",
            user: "AI Agent",
            action: "finished screening",
            target: "Sarah Smith",
            meta: "45% Match",
            time: "1d ago",
            icon: Sparkles,
            color: "text-indigo-500",
            isAi: true,
        },
        {
            type: "job",
            user: "System",
            action: "New Job",
            target: "Product Designer",
            meta: "created",
            time: "1d ago",
            icon: Briefcase,
            color: "text-blue-500",
        },
    ];

    return (
        <Card className="col-span-1 border-white/10 bg-white/5 backdrop-blur-md shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)] hover:shadow-lg transition-all duration-300 dark:bg-white/5 dark:border-white/10">
            <CardHeader>
                <CardTitle>{t("title")}</CardTitle>
            </CardHeader>
            <CardContent>
                <ScrollArea className="h-[350px] pr-4">
                    <div className="space-y-4">
                        {activities.length === 0 ? (
                            <p className="text-sm text-muted-foreground text-center py-4">{t("empty")}</p>
                        ) : (
                            activities.map((activity, index) => (
                                <div
                                    key={index}
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
