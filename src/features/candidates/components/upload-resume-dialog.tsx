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

import { UploadProgress } from "./upload-progress";

interface UploadResumeDialogProps {
    children?: React.ReactNode;
}

export function UploadResumeDialog({ children }: UploadResumeDialogProps) {
    const t = useTranslations("Candidates");
    const [isOpen, setIsOpen] = useState(false);
    const [files, setFiles] = useState<File[]>([]);
    const [isUploading, setIsUploading] = useState(false);

    const onDrop = useCallback((acceptedFiles: File[]) => {
        if (acceptedFiles.length > 0) {
            setFiles(acceptedFiles); // Store all accepted files
            setIsUploading(true);
        }
    }, []);

    const handleUploadComplete = (count: number) => {
        // Wait for user to see the completion state before closing
        // or provide a manual close button. For now auto-close after delay.
        setTimeout(() => {
            setIsOpen(false);
            setTimeout(() => {
                setIsUploading(false);
                setFiles([]);
                toast.success(`Successfully processed ${count} candidates`);
            }, 300);
        }, 1500);
    };

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: { "application/pdf": [".pdf"] },
        multiple: true, // Enable multiple file selection
        disabled: isUploading
    });

    return (
        <Dialog open={isOpen} onOpenChange={(open) => {
            if (!isUploading) setIsOpen(open);
        }}>
            <DialogTrigger asChild>{children}</DialogTrigger>
            <DialogContent className="sm:max-w-md bg-white dark:bg-zinc-950 border border-border shadow-2xl p-8 z-[51] !opacity-100 block">
                <DialogHeader className="mb-6">
                    <DialogTitle>{isUploading ? t("processingTitle", { fallback: "Processing Candidates" }) : t("uploadModalTitle")}</DialogTitle>
                </DialogHeader>

                <AnimatePresence mode="wait">
                    {isUploading ? (
                        <motion.div
                            key="progress"
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, y: -10 }}
                        >
                            <UploadProgress files={files} onComplete={handleUploadComplete} />
                        </motion.div>
                    ) : (
                        <motion.div
                            key="dropzone"
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, y: -10 }}
                        >
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
                        </motion.div>
                    )}
                </AnimatePresence>
            </DialogContent>
        </Dialog>
    );
}
