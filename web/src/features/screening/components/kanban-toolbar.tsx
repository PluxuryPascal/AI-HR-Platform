"use client";

import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";

interface KanbanToolbarProps {
    isSelectionMode: boolean;
    onToggleSelectionMode: () => void;
}

export function KanbanToolbar({ isSelectionMode, onToggleSelectionMode }: KanbanToolbarProps) {
    const t = useTranslations("Screening");

    return (
        <div className="flex justify-end px-2">
            <Button
                variant={isSelectionMode ? "secondary" : "outline"}
                onClick={onToggleSelectionMode}
                size="sm"
            >
                {isSelectionMode ? t("bar.cancelSelection") : t("bar.selectCandidates")}
            </Button>
        </div>
    );
}
