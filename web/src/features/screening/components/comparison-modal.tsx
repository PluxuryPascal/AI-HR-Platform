"use client";

import { useState, useEffect } from "react";
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

interface ComparisonModalProps {
    isOpen: boolean;
    onClose: () => void;
    candidates: CandidateCard[];
}

// Mock AI data generator
const generateMockData = (candidate: CandidateCard) => {
    // Deterministic simulation based on name length/char codes
    const seed = candidate.name.length;

    return {
        skills: ["React", "TypeScript", "Node.js", "Python", "AWS", "Docker"].slice(0, 3 + (seed % 4)),
        experience: `${2 + (seed % 8)} years`,
        risks: seed % 3 === 0 ? "High salary expectations" : seed % 4 === 0 ? "Remote only" : "None identified",
        salary: `$${80 + (seed % 10) * 10}k - $${100 + (seed % 12) * 10}k`,
        summary: seed % 2 === 0
            ? "Strong technical background with proven track record in scalable systems."
            : "Great cultural fit with innovative mindset and rapid learning ability."
    };
};

export function ComparisonModal({ isOpen, onClose, candidates }: ComparisonModalProps) {
    const t = useTranslations("Screening.modal");
    const [isLoading, setIsLoading] = useState(true);
    const [aiData, setAiData] = useState<Record<string, any>>({});
    const { handleExport } = useCsvExport();

    useEffect(() => {
        if (isOpen) {
            setIsLoading(true);
            const timer = setTimeout(() => {
                const data: Record<string, any> = {};
                candidates.forEach(c => {
                    data[c.id] = generateMockData(c);
                });
                setAiData(data);
                setIsLoading(false);
            }, 2000);
            return () => clearTimeout(timer);
        }
    }, [isOpen, candidates]);

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
                            <Button variant="outline" size="sm" onClick={() => handleExport({ candidates, aiData })}>
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
                        <ComparisonTable candidates={candidates} aiData={aiData} />
                        <ScrollBar orientation="horizontal" />
                    </ScrollArea>
                )}
            </DialogContent>
        </Dialog>
    );
}
