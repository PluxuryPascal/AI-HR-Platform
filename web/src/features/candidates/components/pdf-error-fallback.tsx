"use client";

import { Button } from "@/components/ui/button";
import { Download, FileText } from "lucide-react";
import { useTranslations } from "next-intl";

interface PdfErrorFallbackProps {
    url: string;
}

export function PdfErrorFallback({ url }: PdfErrorFallbackProps) {
    const t = useTranslations("Screening.pdfViewer");

    return (
        <div className="w-full h-full flex flex-col items-center justify-center bg-gray-50/50 p-6 text-center space-y-4">
            <div className="bg-white p-4 rounded-full shadow-sm">
                <FileText className="w-8 h-8 text-gray-400" />
            </div>
            <div className="space-y-2">
                <h3 className="font-semibold text-gray-900">{t("error")}</h3>
                <p className="text-sm text-gray-500 max-w-xs mx-auto">
                    The browser couldn&apos;t render the resume directly. You can download it instead.
                </p>
            </div>
            <Button asChild variant="outline">
                <a href={url} download target="_blank" rel="noopener noreferrer">
                    <Download className="w-4 h-4 mr-2" />
                    {t("download")}
                </a>
            </Button>
        </div>
    );
}
