"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { ChevronDown, ChevronUp, RefreshCw, Sparkles } from "lucide-react";

interface InterviewGuideProps {
    matchScore: number;
}

interface Question {
    id: string;
    title: string;
    context: string;
    expected: string;
    type: "expert" | "foundational";
}

export function InterviewGuide({ matchScore }: InterviewGuideProps) {
    const t = useTranslations("InterviewPrep");
    const [loading, setLoading] = useState(true);
    const [questions, setQuestions] = useState<Question[]>([]);
    const [expandedIds, setExpandedIds] = useState<string[]>([]);

    const isExpert = matchScore >= 80;

    const generateQuestions = () => {
        setLoading(true);
        setExpandedIds([]);

        // Simulate AI generation delay
        setTimeout(() => {
            const keys = ["q1", "q2", "q3", "q4", "q5"];
            const type = isExpert ? "expert" : "foundational";

            // In a real app, this would be an API call. 
            // Here we map from translations based on the score.
            const newQuestions = keys.map((key) => ({
                id: `${type}-${key}-${Date.now()}`,
                title: t(`mockQuestions.${type}.${key}.title`),
                context: t(`mockQuestions.${type}.${key}.context`),
                expected: t(`mockQuestions.${type}.${key}.expected`),
                type: type as "expert" | "foundational",
            }));

            setQuestions(newQuestions);
            setLoading(false);
        }, 1500);
    };

    useEffect(() => {
        generateQuestions();
    }, [matchScore]);

    const toggleExpand = (id: string) => {
        setExpandedIds(prev =>
            prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]
        );
    };

    if (loading) {
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
                <Button onClick={generateQuestions} variant="outline" size="sm">
                    <RefreshCw className="w-4 h-4 mr-2" />
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
