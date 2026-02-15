"use client";

import { useEffect, useState, useRef } from "react";
import { Document, Page, pdfjs } from "react-pdf";
import { Button } from "@/components/ui/button";
import { Download, FileText, ChevronLeft, ChevronRight, Loader2, ZoomIn, ZoomOut, RotateCw } from "lucide-react";
import { cn } from "@/lib/utils";
import { useTranslations } from "next-intl";
import { WidgetErrorBoundary } from "@/components/shared/widget-error-boundary";

// Configure worker locally or via CDN (CDN is safer for Next.js to avoid build issues)
pdfjs.GlobalWorkerOptions.workerSrc = `//unpkg.com/pdfjs-dist@${pdfjs.version}/build/pdf.worker.min.mjs`;

interface PDFViewerProps {
    url: string;
}

export const PDFViewer = ({ url }: PDFViewerProps) => {
    const t = useTranslations("Screening.pdfViewer");
    const [numPages, setNumPages] = useState<number>(0);
    const [pageNumber, setPageNumber] = useState<number>(1);
    const [scale, setScale] = useState<number>(1.0);
    const [rotation, setRotation] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<boolean>(false);
    const [containerWidth, setContainerWidth] = useState<number>(0);
    const containerRef = useRef<HTMLDivElement>(null);

    // Responsive width observer with debounce
    useEffect(() => {
        if (!containerRef.current) return;

        let timeoutId: NodeJS.Timeout;

        const resizeObserver = new ResizeObserver((entries) => {
            for (const entry of entries) {
                if (entry.contentRect.width > 0) {
                    // Debounce width updates to prevent rapid re-renders
                    clearTimeout(timeoutId);
                    timeoutId = setTimeout(() => {
                        setContainerWidth(entry.contentRect.width);
                    }, 100);
                }
            }
        });

        resizeObserver.observe(containerRef.current);
        return () => {
            resizeObserver.disconnect();
            clearTimeout(timeoutId);
        };
    }, []);

    const onDocumentLoadSuccess = ({ numPages }: { numPages: number }) => {
        setNumPages(numPages);
        setLoading(false);
        setError(false);
    };

    const onDocumentLoadError = () => {
        setLoading(false);
        setError(true);
    };

    const changePage = (offset: number) => {
        setPageNumber((prevPageNumber) => Math.min(Math.max(1, prevPageNumber + offset), numPages));
    };

    const previousPage = () => changePage(-1);
    const nextPage = () => changePage(1);

    if (error) {
        return (
            <div className="w-full h-full flex flex-col items-center justify-center bg-gray-50/50 p-6 text-center space-y-4">
                <div className="bg-white p-4 rounded-full shadow-sm">
                    <FileText className="w-8 h-8 text-gray-400" />
                </div>
                <div className="space-y-2">
                    <h3 className="font-semibold text-gray-900">{t("error")}</h3>
                    <p className="text-sm text-gray-500 max-w-xs mx-auto">
                        The browser couldn't render the resume directly. You can download it instead.
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

    return (
        <div className="flex flex-col h-full bg-background dark:bg-background relative group transition-colors duration-300">
            {/* Added scrollbar-gutter: stable to prevent layout shift */}
            <div
                className="flex-1 overflow-y-auto overflow-x-hidden relative scrollbar-hide py-8 px-4"
                ref={containerRef}
                style={{ scrollbarGutter: 'stable' }}
            >
                {loading && (
                    <div className="absolute inset-0 flex items-center justify-center z-10 bg-background/50 backdrop-blur-sm transition-all duration-300">
                        <div className="flex flex-col items-center gap-3">
                            <Loader2 className="w-8 h-8 animate-spin text-primary" />
                            <p className="text-sm text-muted-foreground animate-pulse">{t("loading")}</p>
                        </div>
                    </div>
                )}

                <div className={cn("min-h-full flex justify-center", loading ? "opacity-0" : "opacity-100 transition-opacity duration-500")}>
                    {containerWidth > 0 && (
                        <Document
                            file={url}
                            onLoadSuccess={onDocumentLoadSuccess}
                            onLoadError={onDocumentLoadError}
                            className={cn(
                                "rounded-sm overflow-hidden transition-all duration-300",
                                "border border-border/50", // Thin border instead of shadow
                                // Smart Dark Mode Filter - ensure no glowing white outline
                                "dark:filter dark:invert-[0.9] dark:hue-rotate-180 dark:contrast-[.85]"
                            )}
                            loading={null}
                        >
                            <Page
                                pageNumber={pageNumber}
                                // Subtract padding (32px = 2rem) and an extra buffer (2px) to prevent sub-pixel issues
                                width={containerWidth - 34}
                                scale={scale}
                                rotate={rotation}
                                renderTextLayer={false}
                                renderAnnotationLayer={false}
                                className="bg-white"
                            />
                        </Document>
                    )}
                </div>
            </div>

            {/* Floating Toolbar */}
            {!loading && !error && numPages > 0 && (
                <div className="absolute bottom-6 left-1/2 -translate-x-1/2 z-20">
                    <div className="rounded-full p-1.5 flex items-center gap-1 shadow-2xl border border-white/20 dark:border-white/5 backdrop-blur-xl bg-white/70 dark:bg-zinc-900/80 text-foreground transition-all duration-300 hover:scale-105 hover:bg-white/90 dark:hover:bg-zinc-800/90">
                        <div className="flex items-center gap-0.5 px-2 border-r border-border/50 mr-1">
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8 rounded-full hover:bg-black/5 dark:hover:bg-white/10"
                                onClick={previousPage}
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
                                onClick={nextPage}
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
                                onClick={() => setScale(s => Math.max(0.5, s - 0.1))}
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
                                onClick={() => setScale(s => Math.min(2.0, s + 0.1))}
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
                            onClick={() => setRotation(r => (r + 90) % 360)}
                            title={t("rotate")}
                        >
                            <RotateCw className="w-3.5 h-3.5" />
                        </Button>
                    </div>
                </div>
            )}
        </div>
    );
};
