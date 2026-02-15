"use client";

import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useTranslations } from "next-intl";
import { CheckCircle2, FileText, Loader2, AlertCircle } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Progress } from "@/components/ui/progress";

interface FileProgress {
    name: string;
    progress: number;
    status: string;
    isComplete: boolean;
}

interface UseMockUploadProgressProps {
    files: File[];
    onComplete?: (count: number) => void;
}

export function useMockUploadProgress({ files, onComplete }: UseMockUploadProgressProps) {
    const [fileProgress, setFileProgress] = useState<Record<string, FileProgress>>({});
    const t = useTranslations("UploadStages");

    useEffect(() => {
        // Initialize progress state for new files
        const initialProgress: Record<string, FileProgress> = {};
        files.forEach(file => {
            if (!fileProgress[file.name]) {
                initialProgress[file.name] = {
                    name: file.name,
                    progress: 0,
                    status: t("stage1"),
                    isComplete: false
                };
            }
        });

        if (Object.keys(initialProgress).length > 0) {
            setFileProgress(prev => ({ ...prev, ...initialProgress }));
        }

        const intervals: NodeJS.Timeout[] = [];

        files.forEach(file => {
            // Skip if already complete
            if (fileProgress[file.name]?.isComplete) return;

            // Random start delay for realism
            const startDelay = Math.random() * 1000;

            const timeout = setTimeout(() => {
                const interval = setInterval(() => {
                    setFileProgress(prev => {
                        const current = prev[file.name];
                        if (!current || current.isComplete) {
                            clearInterval(interval);
                            return prev;
                        }

                        let increment = 0;
                        if (current.progress < 20) increment = 2 + Math.random();
                        else if (current.progress < 50) increment = 0.5 + Math.random();
                        else if (current.progress < 85) increment = 0.8 + Math.random();
                        else if (current.progress < 100) increment = 2 + Math.random();

                        const newProgress = Math.min(current.progress + increment, 100);

                        let statusKey = "stage1";
                        if (newProgress >= 20 && newProgress < 50) statusKey = "stage2";
                        if (newProgress >= 50 && newProgress < 85) statusKey = "stage3";
                        if (newProgress >= 85) statusKey = "stage4";

                        const isComplete = newProgress >= 100;

                        return {
                            ...prev,
                            [file.name]: {
                                ...current,
                                progress: newProgress,
                                status: t(statusKey),
                                isComplete
                            }
                        };
                    });
                }, 100);
                intervals.push(interval);
            }, startDelay);
        });

        return () => {
            intervals.forEach(clearInterval);
        };
    }, [files, t]); // Removed fileProgress dependency to avoid loops, logic handled inside setters

    const allComplete = files.length > 0 && files.every(f => fileProgress[f.name]?.isComplete);

    useEffect(() => {
        if (allComplete) {
            const completeCount = files.length;
            // Small delay to show 100% state
            const timeout = setTimeout(() => {
                onComplete?.(completeCount);
            }, 800);
            return () => clearTimeout(timeout);
        }
    }, [allComplete, onComplete, files.length]);

    return {
        fileProgress,
        isAllComplete: allComplete
    };
}

interface UploadProgressProps {
    files: File[];
    onComplete: (count: number) => void;
}

export function UploadProgress({ files, onComplete }: UploadProgressProps) {
    const { fileProgress } = useMockUploadProgress({ files, onComplete });

    return (
        <div className="w-full flex flex-col space-y-6">
            <div className="text-center space-y-2">
                <h3 className="text-lg font-medium text-foreground flex items-center justify-center gap-2">
                    {Object.values(fileProgress).every(f => f.isComplete) ? (
                        <CheckCircle2 className="w-5 h-5 text-green-500" />
                    ) : (
                        <Loader2 className="w-5 h-5 text-primary animate-spin" />
                    )}
                    {Object.values(fileProgress).filter(f => f.isComplete).length} / {files.length} Processed
                </h3>
            </div>

            <ScrollArea className="h-[300px] w-full pr-4">
                <div className="space-y-4">
                    <AnimatePresence mode="popLayout">
                        {files.map((file) => {
                            const progress = fileProgress[file.name];
                            if (!progress) return null;

                            return (
                                <motion.div
                                    key={file.name}
                                    initial={{ opacity: 0, y: 10 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    exit={{ opacity: 0, scale: 0.9 }}
                                    className="bg-muted/30 border border-border rounded-lg p-3 space-y-2"
                                >
                                    <div className="flex items-center justify-between gap-3">
                                        <div className="flex items-center gap-3 min-w-0">
                                            <div className={`
                                        h-8 w-8 rounded-full flex items-center justify-center shrink-0 transition-colors
                                        ${progress.isComplete ? "bg-green-100 dark:bg-green-900/30 text-green-600" : "bg-blue-100 dark:bg-blue-900/30 text-blue-600"}
                                    `}>
                                                {progress.isComplete ? <CheckCircle2 className="w-4 h-4" /> : <FileText className="w-4 h-4" />}
                                            </div>
                                            <span className="text-sm font-medium truncate max-w-[200px]">{file.name}</span>
                                        </div>
                                        <span className="text-xs text-muted-foreground whitespace-nowrap tabular-nums">
                                            {Math.round(progress.progress)}%
                                        </span>
                                    </div>

                                    <Progress value={progress.progress} className="h-1.5" />

                                    <p className="text-xs text-muted-foreground w-full text-right animate-pulse">
                                        {progress.status}
                                    </p>
                                </motion.div>
                            );
                        })}
                    </AnimatePresence>
                </div>
            </ScrollArea>
        </div>
    );
}
