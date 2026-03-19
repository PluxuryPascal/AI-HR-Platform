"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { ChevronDown, ChevronUp, RefreshCw, Sparkles } from "lucide-react";
import { useGenerateInterview, useGetInterviewGuide } from "../api/use-candidate-ai";
import { useParams } from "next/navigation";

interface InterviewGuideProps {
    matchScore: number;
    candidateId: string;
}

interface Question {
    id: string;
    title: string;
    context: string;
    expected: string;
    type: string;
}

export function InterviewGuide({ matchScore, candidateId }: InterviewGuideProps) {
    const t = useTranslations("InterviewPrep");
    const [expandedIds, setExpandedIds] = useState<string[]>([]);

    const { data: guide, isLoading } = useGetInterviewGuide(candidateId);
    const { mutate: generate, isPending: isGenerating } = useGenerateInterview();

    const isExpert = matchScore >= 80;

    const handleGenerate = () => {
        generate({ id: candidateId, locale: "ru" });
    };

    const questions: Question[] = guide?.questions.map(q => ({
        id: q.id,
        title: q.question,
        context: q.category,
        expected: q.answer,
        type: q.category.toLowerCase().includes("expert") ? "expert" : "foundational"
    })) || [];

    const toggleExpand = (id: string) => {
        setExpandedIds(prev =>
            prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]
        );
    };

    if (isLoading || isGenerating) {
        return (
            <div className="flex flex-col items-center justify-center p-12 space-y-4 text-center h-full">
                <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
                >
                    <Sparkles className="w-8 h-8 text-primary" />
                </motion.div>
                <div className="space-y-2">
                    <h3 className="text-lg font-medium animate-pulse">{t("loading")}</h3>
                    <p className="text-sm text-muted-foreground">{t("loadingDesc")}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-6 p-4">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-xl font-semibold flex items-center gap-2">
                        {t("title")}
                        <Badge variant={isExpert ? "default" : "secondary"}>
                            {isExpert ? t("expertLevel") : t("foundationalLevel")}
                        </Badge>
                    </h2>
                    <p className="text-sm text-muted-foreground mt-1">
                        {t("description")}
                    </p>
                </div>
                <Button onClick={handleGenerate} variant="outline" size="sm" disabled={isGenerating}>
                    {isGenerating ? <RefreshCw className="w-4 h-4 mr-2 animate-spin" /> : <RefreshCw className="w-4 h-4 mr-2" />}
                    {t("regenerate")}
                </Button>
            </div>

            <div className="grid gap-4">
                {questions.map((q, index) => (
                    <motion.div
                        key={q.id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                    >
                        <Card className="overflow-hidden border-l-4 border-l-primary/50">
                            <div
                                className="p-4 flex items-start justify-between cursor-pointer hover:bg-muted/50 transition-colors"
                                onClick={() => toggleExpand(q.id)}
                            >
                                <div className="space-y-1">
                                    <div className="flex items-center gap-2">
                                        <Badge variant="outline" className="text-[10px] uppercase tracking-wider">
                                            {t("card.context")}
                                        </Badge>
                                        <span className="text-xs text-muted-foreground">{q.context}</span>
                                    </div>
                                    <h3 className="font-medium text-base">{q.title}</h3>
                                </div>
                                <Button variant="ghost" size="sm" className="shrink-0 ml-2">
                                    {expandedIds.includes(q.id) ? (
                                        <ChevronUp className="w-4 h-4" />
                                    ) : (
                                        <ChevronDown className="w-4 h-4" />
                                    )}
                                </Button>
                            </div>

                            <AnimatePresence>
                                {expandedIds.includes(q.id) && (
                                    <motion.div
                                        initial={{ height: 0, opacity: 0 }}
                                        animate={{ height: "auto", opacity: 1 }}
                                        exit={{ height: 0, opacity: 0 }}
                                        transition={{ duration: 0.2 }}
                                    >
                                        <div className="px-4 pb-4 pt-0 border-t bg-muted/20">
                                            <div className="mt-4 space-y-2">
                                                <h4 className="text-sm font-semibold text-primary flex items-center gap-2">
                                                    <Sparkles className="w-3 h-3" />
                                                    {t("card.expected")}
                                                </h4>
                                                <p className="text-sm text-muted-foreground leading-relaxed">
                                                    {q.expected}
                                                </p>
                                            </div>
                                        </div>
                                    </motion.div>
                                )}
                            </AnimatePresence>
                        </Card>
                    </motion.div>
                ))}
            </div>
        </div>
    );
}
