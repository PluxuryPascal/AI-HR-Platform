"use client";

import { useState, useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { useTranslations } from "next-intl";
import { motion, AnimatePresence } from "framer-motion";
import { Upload, FileText, CheckCircle, X } from "lucide-react";
import { toast } from "sonner";

import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface UploadResumeDialogProps {
    children?: React.ReactNode;
}

interface UploadingFile {
    name: string;
    progress: number;
    status: "uploading" | "completed";
}

export function UploadResumeDialog({ children }: UploadResumeDialogProps) {
    const t = useTranslations("Candidates");
    const [isOpen, setIsOpen] = useState(false);
    const [files, setFiles] = useState<UploadingFile[]>([]);

    const onDrop = useCallback((acceptedFiles: File[]) => {
        const newFiles = acceptedFiles.map((file) => ({
            name: file.name,
            progress: 0,
            status: "uploading" as const,
        }));

        setFiles((prev) => [...prev, ...newFiles]);

        // Simulate upload progress
        newFiles.forEach((file, index) => {
            const duration = 2000; // 2 seconds
            const interval = 50;
            const steps = duration / interval;
            let step = 0;

            const timer = setInterval(() => {
                step++;
                const progress = Math.min((step / steps) * 100, 100);

                setFiles((prev) =>
                    prev.map((f) =>
                        f.name === file.name
                            ? {
                                ...f,
                                progress,
                                status: progress === 100 ? "completed" : "uploading",
                            }
                            : f
                    )
                );

                if (step >= steps) {
                    clearInterval(timer);
                    if (index === newFiles.length - 1) {
                        toast.success("Resumes uploaded successfully");
                    }
                }
            }, interval);
        });
    }, []);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: { "application/pdf": [".pdf"] },
    });

    return (
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>{children}</DialogTrigger>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>{t("uploadModalTitle")}</DialogTitle>
                </DialogHeader>

                <div
                    {...getRootProps()}
                    className={`
            border-2 border-dashed rounded-lg p-10 text-center cursor-pointer transition-colors
            ${isDragActive
                            ? "border-blue-500 bg-blue-50/50 dark:bg-blue-900/20"
                            : "border-muted-foreground/25 hover:border-blue-500/50 hover:bg-muted/50"
                        }
          `}
                >
                    <input {...getInputProps()} />
                    <div className="flex flex-col items-center justify-center gap-2 text-muted-foreground">
                        <Upload className="h-10 w-10 text-muted-foreground/50" />
                        <p className="text-sm font-medium">{t("dropzoneText")}</p>
                    </div>
                </div>

                <div className="space-y-3 mt-4 max-h-[200px] overflow-y-auto pr-2">
                    <AnimatePresence>
                        {files.map((file) => (
                            <motion.div
                                key={file.name}
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, x: -10 }}
                                className="flex items-center gap-3 p-3 rounded-md bg-muted/30 border border-border"
                            >
                                <div className="h-8 w-8 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center shrink-0">
                                    <FileText className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="flex justify-between items-center mb-1">
                                        <p className="text-sm font-medium truncate">{file.name}</p>
                                        {file.status === "completed" ? (
                                            <CheckCircle className="h-4 w-4 text-green-500" />
                                        ) : (
                                            <span className="text-xs text-muted-foreground">
                                                {Math.round(file.progress)}%
                                            </span>
                                        )}
                                    </div>
                                    <div className="h-1.5 w-full bg-muted rounded-full overflow-hidden">
                                        <motion.div
                                            className="h-full bg-blue-600 rounded-full"
                                            initial={{ width: 0 }}
                                            animate={{ width: `${file.progress}%` }}
                                            transition={{ duration: 0.1 }}
                                        />
                                    </div>
                                </div>
                            </motion.div>
                        ))}
                    </AnimatePresence>
                </div>
            </DialogContent>
        </Dialog>
    );
}
