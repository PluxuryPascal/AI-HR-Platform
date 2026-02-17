"use client";

import { useTranslations } from "next-intl";
import { X, Trash2, ArrowRight } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface BulkActionBarProps {
    selectedCount: number;
    onClear: () => void;
    onRejectAll: () => void;
    onMoveTo: (status: string) => void;
}

export function BulkActionBar({ selectedCount, onClear, onRejectAll, onMoveTo }: BulkActionBarProps) {
    const t = useTranslations("BulkActions");

    return (
        <AnimatePresence>
            {selectedCount > 0 && (
                <motion.div
                    initial={{ y: 100, opacity: 0, x: "-50%" }}
                    animate={{ y: 0, opacity: 1, x: "-50%" }}
                    exit={{ y: 100, opacity: 0, x: "-50%" }}
                    className="fixed bottom-8 left-1/2 z-50 flex items-center gap-4 rounded-full border border-border bg-popover px-6 py-3 text-popover-foreground shadow-2xl"
                >
                    <div className="flex items-center gap-2 border-r pr-4">
                        <span className="font-medium whitespace-nowrap">
                            {t("selectedCount", { count: selectedCount })}
                        </span>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6 rounded-full hover:bg-muted"
                            onClick={onClear}
                        >
                            <X className="h-4 w-4" />
                            <span className="sr-only">{t("clear")}</span>
                        </Button>
                    </div>

                    <div className="flex items-center gap-2">
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="outline" size="sm" className="rounded-full">
                                    {t("moveTo")} <ArrowRight className="ml-2 h-4 w-4" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuItem onClick={() => onMoveTo("screening")}>
                                    Screening
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => onMoveTo("interview")}>
                                    Interview
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => onMoveTo("offer")}>
                                    Offer
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>

                        <Button
                            variant="destructive"
                            size="sm"
                            className="rounded-full"
                            onClick={onRejectAll}
                        >
                            <Trash2 className="mr-2 h-4 w-4" />
                            {t("rejectAll")}
                        </Button>
                    </div>
                </motion.div>
            )}
        </AnimatePresence>
    );
}
