"use client";

import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight, ZoomIn, ZoomOut, RotateCw } from "lucide-react";
import { useTranslations } from "next-intl";
import { glassToolbar } from "@/lib/styles";

interface PdfToolbarProps {
    pageNumber: number;
    numPages: number;
    scale: number;
    onPreviousPage: () => void;
    onNextPage: () => void;
    onZoomIn: () => void;
    onZoomOut: () => void;
    onRotate: () => void;
}

export function PdfToolbar({
    pageNumber,
    numPages,
    scale,
    onPreviousPage,
    onNextPage,
    onZoomIn,
    onZoomOut,
    onRotate,
}: PdfToolbarProps) {
    const t = useTranslations("Screening.pdfViewer");

    return (
        <div className="absolute bottom-6 left-1/2 -translate-x-1/2 z-20">
            <div className={`rounded-full p-1.5 flex items-center gap-1 ${glassToolbar} hover:scale-105 hover:bg-white/90 dark:hover:bg-zinc-800/90`}>
                <div className="flex items-center gap-0.5 px-2 border-r border-border/50 mr-1">
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 rounded-full hover:bg-black/5 dark:hover:bg-white/10"
                        onClick={onPreviousPage}
                        disabled={pageNumber <= 1}
                    >
                        <ChevronLeft className="w-4 h-4" />
                    </Button>
                    <span className="text-xs font-medium tabular-nums px-2 min-w-[3rem] text-center">
                        {t("page")} {pageNumber} {t("of")} {numPages}
                    </span>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 rounded-full hover:bg-black/5 dark:hover:bg-white/10"
                        onClick={onNextPage}
                        disabled={pageNumber >= numPages}
                    >
                        <ChevronRight className="w-4 h-4" />
                    </Button>
                </div>

                <div className="flex items-center gap-0.5">
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 rounded-full hover:bg-black/5 dark:hover:bg-white/10"
                        onClick={onZoomOut}
                        title={t("zoomOut")}
                    >
                        <ZoomOut className="w-3.5 h-3.5" />
                    </Button>
                    <span className="text-xs font-medium tabular-nums w-12 text-center">
                        {Math.round(scale * 100)}%
                    </span>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 rounded-full hover:bg-black/5 dark:hover:bg-white/10"
                        onClick={onZoomIn}
                        title={t("zoomIn")}
                    >
                        <ZoomIn className="w-3.5 h-3.5" />
                    </Button>
                </div>

                <div className="w-px h-4 bg-border/50 mx-1" />

                <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 rounded-full hover:bg-black/5 dark:hover:bg-white/10"
                    onClick={onRotate}
                    title={t("rotate")}
                >
                    <RotateCw className="w-3.5 h-3.5" />
                </Button>
            </div>
        </div>
    );
}
