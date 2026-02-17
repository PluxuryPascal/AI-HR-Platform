"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { ScoreFactor } from "@/features/screening/types";
import { useTranslations } from "next-intl";
import { Zap, ArrowRight } from "lucide-react";
import { Link } from "@/i18n/routing";

interface MatchScoreBadgeProps {
    score: number;
    breakdown?: ScoreFactor[];
    className?: string;
    candidateId?: string; // Optional: for linking to details
}

export function MatchScoreBadge({ score, breakdown, className, candidateId }: MatchScoreBadgeProps) {
    const t = useTranslations("MatchScore");

    let colorClass = "bg-red-500/15 text-red-700 hover:bg-red-500/25 border-red-200 dark:bg-red-500/20 dark:text-red-400 dark:border-red-800";

    if (score >= 80) {
        colorClass = "bg-green-500/15 text-green-700 hover:bg-green-500/25 border-green-200 dark:bg-green-500/20 dark:text-green-400 dark:border-green-800";
    } else if (score >= 50) {
        colorClass = "bg-yellow-500/15 text-yellow-700 hover:bg-yellow-500/25 border-yellow-200 dark:bg-yellow-500/20 dark:text-yellow-400 dark:border-yellow-800";
    }

    if (!breakdown || breakdown.length === 0) {
        return (
            <Badge variant="outline" className={cn("font-medium", colorClass, className)}>
                {score}% Match
            </Badge>
        );
    }

    return (
        <Popover>
            <PopoverTrigger asChild>
                <Badge
                    variant="outline"
                    className={cn("font-medium cursor-pointer transition-colors", colorClass, className)}
                >
                    {score}% Match
                </Badge>
            </PopoverTrigger>
            <PopoverContent
                className="w-80 sm:w-96 p-4 bg-popover border-border shadow-2xl z-50"
                side="bottom"
                align="center"
            >
                <div className="space-y-4">
                    {/* Header */}
                    <div className="flex items-center gap-2 pb-1">
                        <div className="p-1.5 bg-primary/10 rounded-full">
                            <Zap className="w-4 h-4 text-primary" />
                        </div>
                        <h4 className="font-semibold text-sm">{t("analysisTitle")}</h4>
                    </div>

                    {/* Breakdown List */}
                    <div className="space-y-3">
                        {breakdown.map((factor) => (
                            <div key={factor.id} className="flex items-start justify-between gap-4">
                                <span className="flex-1 text-sm text-balance leading-tight text-muted-foreground">
                                    {factor.text}
                                </span>
                                <span className={cn(
                                    "shrink-0 text-xs font-medium px-2 py-1 rounded-md",
                                    factor.type === 'positive' && "text-emerald-700 bg-emerald-100 dark:text-emerald-400 dark:bg-emerald-950/60",
                                    factor.type === 'negative' && "text-rose-700 bg-rose-100 dark:text-rose-400 dark:bg-rose-950/60",
                                    factor.type === 'neutral' && "text-gray-700 bg-gray-100 dark:text-gray-300 dark:bg-gray-800"
                                )}>
                                    {factor.impact}
                                </span>
                            </div>
                        ))}
                    </div>

                    <Separator />

                    {/* Footer Action */}
                    <div className="flex justify-end">
                        {candidateId ? (
                            <Button variant="ghost" size="sm" className="h-auto p-0 text-primary hover:text-primary/80 hover:bg-transparent font-normal group" asChild>
                                <Link href={`/dashboard/candidates/${candidateId}`}>
                                    {t("viewFullAnalysis")}
                                    <ArrowRight className="w-3 h-3 ml-1 transition-transform group-hover:translate-x-0.5" />
                                </Link>
                            </Button>
                        ) : (
                            <Button variant="ghost" size="sm" className="h-auto p-0 text-primary hover:text-primary/80 hover:bg-transparent font-normal group">
                                {t("viewFullAnalysis")}
                                <ArrowRight className="w-3 h-3 ml-1 transition-transform group-hover:translate-x-0.5" />
                            </Button>
                        )}
                    </div>
                </div>
            </PopoverContent>
        </Popover>
    );
}
