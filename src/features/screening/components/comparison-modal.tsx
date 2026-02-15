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
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Loader2, Download, Sparkles } from "lucide-react";
import { CandidateCard } from "../types";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

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
    const tCriteria = useTranslations("Screening.criteria");
    const [isLoading, setIsLoading] = useState(true);
    const [aiData, setAiData] = useState<Record<string, any>>({});

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

    const handleExport = () => {
        const headers = ["Criteria", ...candidates.map(c => c.name)];
        const rows = [
            ["Role", ...candidates.map(c => c.role)],
            ["Score", ...candidates.map(c => c.score.toString())],
            ["Experience", ...candidates.map(c => aiData[c.id]?.experience || "")],
            ["Skills", ...candidates.map(c => aiData[c.id]?.skills.join(", ") || "")],
            ["Salary", ...candidates.map(c => aiData[c.id]?.salary || "")],
            ["Risks", ...candidates.map(c => aiData[c.id]?.risks || "")],
            ["Summary", ...candidates.map(c => aiData[c.id]?.summary || "")],
        ];

        const csvContent = [
            headers.join(","),
            ...rows.map(r => r.map(cell => `"${cell}"`).join(","))
        ].join("\n");

        const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);
        link.setAttribute("href", url);
        link.setAttribute("download", "candidate_comparison.csv");
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

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
                            <Button variant="outline" size="sm" onClick={handleExport}>
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
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[200px]">Criteria</TableHead>
                                    {candidates.map(candidate => (
                                        <TableHead key={candidate.id} className="min-w-[250px]">
                                            <div className="flex items-center gap-3">
                                                <Avatar>
                                                    <AvatarImage src={candidate.avatarUrl} />
                                                    <AvatarFallback>{candidate.name[0]}</AvatarFallback>
                                                </Avatar>
                                                <div className="flex flex-col">
                                                    <span className="font-medium text-foreground">{candidate.name}</span>
                                                    <span className="text-xs font-normal">{candidate.role}</span>
                                                </div>
                                            </div>
                                        </TableHead>
                                    ))}
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                <TableRow>
                                    <TableCell className="font-medium">Score</TableCell>
                                    {candidates.map(c => (
                                        <TableCell key={c.id}>
                                            <Badge variant={c.score >= 80 ? "default" : "secondary"}>
                                                {c.score}
                                            </Badge>
                                        </TableCell>
                                    ))}
                                </TableRow>
                                <TableRow>
                                    <TableCell className="font-medium">{tCriteria("experience")}</TableCell>
                                    {candidates.map(c => (
                                        <TableCell key={c.id}>{aiData[c.id]?.experience}</TableCell>
                                    ))}
                                </TableRow>
                                <TableRow>
                                    <TableCell className="font-medium">{tCriteria("skills")}</TableCell>
                                    {candidates.map(c => (
                                        <TableCell key={c.id}>
                                            <div className="flex flex-wrap gap-1">
                                                {aiData[c.id]?.skills.map((skill: string) => (
                                                    <Badge key={skill} variant="outline" className="text-xs">
                                                        {skill}
                                                    </Badge>
                                                ))}
                                            </div>
                                        </TableCell>
                                    ))}
                                </TableRow>
                                <TableRow>
                                    <TableCell className="font-medium">{tCriteria("salary")}</TableCell>
                                    {candidates.map(c => (
                                        <TableCell key={c.id}>{aiData[c.id]?.salary}</TableCell>
                                    ))}
                                </TableRow>
                                <TableRow>
                                    <TableCell className="font-medium text-destructive">{tCriteria("risks")}</TableCell>
                                    {candidates.map(c => (
                                        <TableCell key={c.id} className="text-destructive/80">
                                            {aiData[c.id]?.risks}
                                        </TableCell>
                                    ))}
                                </TableRow>
                                <TableRow>
                                    <TableCell className="font-medium text-primary">
                                        <div className="flex items-center gap-2">
                                            <Sparkles className="w-3 h-3" />
                                            {tCriteria("summary")}
                                        </div>
                                    </TableCell>
                                    {candidates.map(c => (
                                        <TableCell key={c.id} className="italic text-muted-foreground">
                                            "{aiData[c.id]?.summary}"
                                        </TableCell>
                                    ))}
                                </TableRow>
                            </TableBody>
                        </Table>
                        <ScrollBar orientation="horizontal" />
                    </ScrollArea>
                )}
            </DialogContent>
        </Dialog>
    );
}
