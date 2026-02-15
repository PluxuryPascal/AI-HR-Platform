"use client";

import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";

interface FloatingComparisonBarProps {
    selectedCount: number;
    onCompare: () => void;
    onClear: () => void;
}

export function FloatingComparisonBar({
    selectedCount,
    onCompare,
    onClear
}: FloatingComparisonBarProps) {
    const t = useTranslations("Screening.bar");

    if (selectedCount === 0) return null;

    return (
        <div className="fixed bottom-6 left-1/2 transform -translate-x-1/2 bg-background border shadow-lg rounded-full px-6 py-3 flex items-center space-x-4 z-50 animate-in slide-in-from-bottom-5 fade-in">
            <span className="text-sm font-medium">
                {t("selected", { count: selectedCount })}
            </span>
            <div className="flex items-center gap-2">
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onClear}
                    className="rounded-full px-3"
                >
                    {t("clear")}
                </Button>
                <Button
                    onClick={onCompare}
                    disabled={selectedCount < 2 || selectedCount > 10}
                    className="rounded-full"
                >
                    {t("compareBtn")}
                </Button>
            </div>
        </div>
    );
}
