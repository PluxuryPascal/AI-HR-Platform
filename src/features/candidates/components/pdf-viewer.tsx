import { Button } from "@/components/ui/button";
import { Download, FileText } from "lucide-react";
import { useState } from "react";

interface PDFViewerProps {
    url: string;
}

export const PDFViewer = ({ url }: PDFViewerProps) => {
    const [error, setError] = useState(false);

    if (error) {
        return (
            <div className="w-full h-full flex flex-col items-center justify-center bg-gray-50/50 p-6 text-center space-y-4">
                <div className="bg-white p-4 rounded-full shadow-sm">
                    <FileText className="w-8 h-8 text-gray-400" />
                </div>
                <div className="space-y-2">
                    <h3 className="font-semibold text-gray-900">Unable to load PDF</h3>
                    <p className="text-sm text-gray-500 max-w-xs mx-auto">
                        The browser couldn't render the resume directly. You can download it instead.
                    </p>
                </div>
                <Button asChild variant="outline">
                    <a href={url} download target="_blank" rel="noopener noreferrer">
                        <Download className="w-4 h-4 mr-2" />
                        Download Resume
                    </a>
                </Button>
            </div>
        );
    }

    return (
        <iframe
            src={url}
            className="w-full h-full border-none bg-white dark:bg-gray-100 filter dark:invert-[0.05] dark:contrast-[.9]"
            title="Resume PDF"
            onError={() => setError(true)}
        />
    );
};
