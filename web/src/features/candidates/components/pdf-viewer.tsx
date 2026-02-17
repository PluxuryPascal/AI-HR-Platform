"use client";

import { Document, Page, pdfjs } from "react-pdf";
import { Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { useTranslations } from "next-intl";
import { PdfToolbar } from "./pdf-toolbar";
import { PdfErrorFallback } from "./pdf-error-fallback";
import { usePdfControls } from "../hooks/use-pdf-controls";

// Configure worker locally or via CDN (CDN is safer for Next.js to avoid build issues)
pdfjs.GlobalWorkerOptions.workerSrc = `//unpkg.com/pdfjs-dist@${pdfjs.version}/build/pdf.worker.min.mjs`;

interface PDFViewerProps {
    url: string;
}

export const PDFViewer = ({ url }: PDFViewerProps) => {
    const t = useTranslations("Screening.pdfViewer");
    const {
        numPages,
        pageNumber,
        scale,
        rotation,
        loading,
        error,
        containerWidth,
        containerRef,
        onDocumentLoadSuccess,
        onDocumentLoadError,
        previousPage,
        nextPage,
        zoomIn,
        zoomOut,
        rotate,
    } = usePdfControls();

    if (error) {
        return <PdfErrorFallback url={url} />;
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
                <PdfToolbar
                    pageNumber={pageNumber}
                    numPages={numPages}
                    scale={scale}
                    onPreviousPage={previousPage}
                    onNextPage={nextPage}
                    onZoomIn={zoomIn}
                    onZoomOut={zoomOut}
                    onRotate={rotate}
                />
            )}
        </div>
    );
};
