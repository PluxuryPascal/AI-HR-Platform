"use client";

import { ChevronLeft, ChevronRight, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
    const t = useTranslations("Common.pagination");

    if (totalPages <= 1) return null;

    const renderPageButtons = () => {
        const buttons = [];
        const maxVisible = 5;

        let start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
        let end = Math.min(totalPages, start + maxVisible - 1);

        if (end - start + 1 < maxVisible) {
            start = Math.max(1, end - maxVisible + 1);
        }

        for (let i = start; i <= end; i++) {
            buttons.push(
                <Button
                    key={i}
                    variant={currentPage === i ? "default" : "outline"}
                    size="sm"
                    onClick={() => onPageChange(i)}
                    className="w-10 h-10"
                >
                    {i}
                </Button>
            );
        }
        return buttons;
    };

    return (
        <div className="flex items-center justify-between px-2 py-4">
            <div className="text-sm text-muted-foreground">
                {t("pageOf", { current: currentPage, total: totalPages })}
            </div>
            <div className="flex items-center space-x-2">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onPageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                >
                    <ChevronLeft className="h-4 w-4 mr-2" />
                    {t("previous")}
                </Button>
                <div className="flex items-center space-x-1">
                    {renderPageButtons()}
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onPageChange(currentPage + 1)}
                    disabled={currentPage === totalPages}
                >
                    {t("next")}
                    <ChevronRight className="h-4 w-4 ml-2" />
                </Button>
            </div>
        </div>
    );
}
