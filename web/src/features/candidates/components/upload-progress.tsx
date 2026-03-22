"use client";

import { useEffect, useState, useRef } from "react";
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

import { useUploadResume } from "../api/use-upload-resume";

interface UploadProgressProps {
    files: File[];
    jobId: string;
    onComplete: (count: number) => void;
}

export function useUploadProgress({ files, jobId, onComplete }: UploadProgressProps) {
    const [fileProgress, setFileProgress] = useState<Record<string, FileProgress>>({});
    const uploadMutation = useUploadResume();
    const t = useTranslations("UploadStages");
    const startedUploads = useRef<Set<string>>(new Set());

    useEffect(() => {
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

        files.forEach(async (file) => {
            // Prevent double upload of the same file in this lifecycle
            if (startedUploads.current.has(file.name)) return;
            startedUploads.current.add(file.name);
            // Simulate initial progress
            let progress = 0;
            const progressInterval = setInterval(() => {
                setFileProgress(prev => {
                    const curr = prev[file.name];
                    if (!curr || curr.isComplete || curr.progress >= 90) {
                        clearInterval(progressInterval);
                        return prev;
                    }
                    const nextProgress = Math.min(curr.progress + Math.random() * 5, 90);
                    return {
                        ...prev,
                        [file.name]: { ...curr, progress: nextProgress, status: t("stage2") }
                    };
                });
            }, 500);

            try {
                await uploadMutation.mutateAsync({ jobId, file });
                setFileProgress(prev => ({
                    ...prev,
                    [file.name]: {
                        ...prev[file.name],
                        progress: 100,
                        status: t("stage4"),
                        isComplete: true
                    }
                }));
            } catch (err) {
                clearInterval(progressInterval);
                setFileProgress(prev => ({
                    ...prev,
                    [file.name]: {
                        ...prev[file.name],
                        status: "Error",
                        isComplete: false
                    }
                }));
            }
        });
    }, [files, jobId]);

    const allComplete = files.length > 0 && Object.values(fileProgress).every(f => f.isComplete);

    useEffect(() => {
        if (allComplete) {
            const timeout = setTimeout(() => {
                onComplete?.(files.length);
            }, 800);
            return () => clearTimeout(timeout);
        }
    }, [allComplete, onComplete, files.length]);

    return { fileProgress };
}


export function UploadProgress({ files, jobId, onComplete }: UploadProgressProps) {
    const { fileProgress } = useUploadProgress({ files, jobId, onComplete });

    const progressValues = Object.values(fileProgress) as FileProgress[];
    const completeCount = progressValues.filter(f => f.isComplete).length;
    const isAllComplete = progressValues.length > 0 && progressValues.every(f => f.isComplete);

    return (
        <div className="w-full flex flex-col space-y-6">
            <div className="text-center space-y-2">
                <h3 className="text-lg font-medium text-foreground flex items-center justify-center gap-2">
                    {isAllComplete ? (
                        <CheckCircle2 className="w-5 h-5 text-green-500" />
                    ) : (
                        <Loader2 className="w-5 h-5 text-primary animate-spin" />
                    )}
                    {completeCount} / {files.length} Processed
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
