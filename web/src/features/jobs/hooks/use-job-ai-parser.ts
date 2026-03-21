import { useState, useCallback } from "react";
import { UseFormReturn } from "react-hook-form";
import { JobFormValues } from "@/features/jobs/schemas/job-schema";
import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types";

export function useJobAiParser(form: UseFormReturn<JobFormValues>) {
    const [isAnalyzing, setIsAnalyzing] = useState(false);
    const [rawDescription, setRawDescription] = useState("");

    const onAnalyze = useCallback(async () => {
        if (!rawDescription.trim()) {
            return { success: false, reason: "empty" as const };
        }

        setIsAnalyzing(true);
        try {
            const response = await apiClient.post<ApiResponse<{
                title: string;
                description: string;
                requirements: string[];
                work_format: string;
            }>>("/ai/parse-job", {
                raw_text: rawDescription,
                locale: "ru" // TODO: get from next-intl if needed
            });

            const result = response.data;

            // Auto-fill fields
            if (result.title) form.setValue("title", result.title);
            if (result.description) form.setValue("description", result.description);
            if (result.requirements) form.setValue("requirements", result.requirements);
            
            // Map work format from backend string to frontend enum
            if (result.work_format) {
                const formatMap: Record<string, "Remote" | "Onsite" | "Hybrid"> = {
                    "remote": "Remote",
                    "office": "Onsite",
                    "hybrid": "Hybrid"
                };
                const mappedFormat = formatMap[result.work_format.toLowerCase()];
                if (mappedFormat) {
                    form.setValue("type", mappedFormat);
                }
            }

            return { success: true as const };
        } catch (error) {
            console.error("AI Parse Error:", error);
            return { success: false, reason: "error" as const };
        } finally {
            setIsAnalyzing(false);
        }
    }, [rawDescription, form]);

    return {
        isAnalyzing,
        rawDescription,
        setRawDescription,
        onAnalyze,
    };
}
