"use client";

import { useEffect } from "react";
import { useTranslations } from "next-intl";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Loader2, Download, Sparkles } from "lucide-react";
import { CandidateCard } from "../types";
import { ComparisonTable } from "./comparison-table";
import { useCsvExport } from "../hooks/use-csv-export";
import { useCompareCandidates } from "../../candidates/api/use-compare-candidates";

interface ComparisonModalProps {
    isOpen: boolean;
    onClose: () => void;
    candidates: CandidateCard[];
    jobId: string;
}


export function ComparisonModal({ isOpen, onClose, candidates, jobId }: ComparisonModalProps) {
    const t = useTranslations("Screening.modal");
    const { handleExport } = useCsvExport();
    
    const { mutate: compare, data: aiData, isPending: isLoading } = useCompareCandidates(jobId);

    useEffect(() => {
        if (isOpen && candidates.length >= 2) {
            compare({
                candidate_ids: candidates.map(c => c.id)
            });
        }
    }, [isOpen, candidates, compare]);

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="max-w-[95vw] w-[95vw] h-[85vh] flex flex-col sm:max-w-[95vw]">
                <DialogHeader>
                    <div className="flex items-center justify-between">
                        <DialogTitle className="flex items-center gap-2">
                            <Sparkles className="w-5 h-5 text-primary" />
                            {t("title")}
                        </DialogTitle>
                        {!isLoading && (
                            <Button variant="outline" size="sm" onClick={() => handleExport({ candidates, aiData: aiData || {} })}>
                                <Download className="w-4 h-4 mr-2" />
                                {t("exportBtn")}
                            </Button>
                        )}
                    </div>
                    <DialogDescription>
                        {isLoading ? t("loading") : t("description", { count: candidates.length })}
                    </DialogDescription>
                </DialogHeader>

                {isLoading ? (
                    <div className="flex-1 flex flex-col items-center justify-center space-y-4">
                        <Loader2 className="w-12 h-12 animate-spin text-primary" />
                        <p className="text-muted-foreground animate-pulse">{t("loading")}</p>
                    </div>
                ) : (
                    <ScrollArea className="flex-1 -mx-6 px-6">
                        <ComparisonTable candidates={candidates} aiData={aiData || {}} />
                        <ScrollBar orientation="horizontal" />
                    </ScrollArea>
                )}
            </DialogContent>
        </Dialog>
    );
}
