"use client";

import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import { Settings2, Plus, Check } from "lucide-react";

interface KanbanToolbarProps {
    isSelectionMode: boolean;
    onToggleSelectionMode: () => void;
    isEditMode: boolean;
    onToggleEditMode: () => void;
    onAddStage: () => void;
}

export function KanbanToolbar({ 
    isSelectionMode, 
    onToggleSelectionMode,
    isEditMode,
    onToggleEditMode,
    onAddStage
}: KanbanToolbarProps) {
    const t = useTranslations("Screening");

    return (
        <div className="flex justify-end px-2 gap-2">
            <Button
                variant={isEditMode ? "secondary" : "outline"}
                onClick={onToggleEditMode}
                size="sm"
                className="gap-2"
            >
                {isEditMode ? (
                    <>
                        <Check className="w-4 h-4" />
                        {t("exitEditMode")}
                    </>
                ) : (
                    <>
                        <Settings2 className="w-4 h-4" />
                        {t("editMode")}
                    </>
                )}
            </Button>
            
            {isEditMode && (
                <Button onClick={onAddStage} size="sm" className="gap-2">
                    <Plus className="w-4 h-4" />
                    {t("addStage")}
                </Button>
            )}

            {!isEditMode && (
                <Button
                    variant={isSelectionMode ? "secondary" : "outline"}
                    onClick={onToggleSelectionMode}
                    size="sm"
                >
                    {isSelectionMode ? t("bar.cancelSelection") : t("bar.selectCandidates")}
                </Button>
            )}
        </div>
    );
}
