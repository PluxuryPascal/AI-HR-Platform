"use client";

import { useTranslations } from "next-intl";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Briefcase, UserCheck, MessageSquare, Clock } from "lucide-react";

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
            icon: UserCheck,
            color: "text-green-500",
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
            icon: UserCheck,
            color: "text-green-500",
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
        <Card className="col-span-1">
            <CardHeader>
                <CardTitle>{t("title")}</CardTitle>
            </CardHeader>
            <CardContent>
                <ScrollArea className="h-[350px] pr-4">
                    <div className="space-y-6">
                        {activities.length === 0 ? (
                            <p className="text-sm text-muted-foreground text-center py-4">{t("empty")}</p>
                        ) : (
                            activities.map((activity, index) => (
                                <div key={index} className="flex items-start gap-4">
                                    <div className={`mt-0.5 rounded-full p-2 bg-muted ${activity.color?.replace("text-", "bg-opacity-10 bg-")}`}>
                                        <activity.icon className={`h-4 w-4 ${activity.color}`} />
                                    </div>
                                    <div className="space-y-1">
                                        <p className="text-sm font-medium leading-none">
                                            <span className="font-semibold">{activity.user}</span>{" "}
                                            <span className="text-muted-foreground">{activity.action}</span>{" "}
                                            <span className="font-semibold">{activity.target}</span>
                                            {activity.meta && (
                                                <span className="text-muted-foreground"> ({activity.meta})</span>
                                            )}
                                        </p>
                                        <div className="flex items-center pt-1 text-xs text-muted-foreground">
                                            <Clock className="mr-1 h-3 w-3" />
                                            {activity.time}
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
