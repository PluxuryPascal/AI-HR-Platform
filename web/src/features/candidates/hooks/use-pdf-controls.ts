import { useState, useEffect, useRef } from "react";

export function usePdfControls() {
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

    const zoomIn = () => setScale(s => Math.min(2.0, s + 0.1));
    const zoomOut = () => setScale(s => Math.max(0.5, s - 0.1));
    const rotate = () => setRotation(r => (r + 90) % 360);

    return {
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
    };
}
