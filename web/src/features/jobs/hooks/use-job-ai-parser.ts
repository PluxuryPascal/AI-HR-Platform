import { useState, useCallback } from "react";
import { UseFormReturn } from "react-hook-form";
import { parseJobDescription } from "@/features/jobs/utils/mock-ai-parser";
import { JobFormValues } from "@/features/jobs/schemas/job-schema";

export function useJobAiParser(form: UseFormReturn<JobFormValues>) {
    const [isAnalyzing, setIsAnalyzing] = useState(false);
    const [rawDescription, setRawDescription] = useState("");

    const onAnalyze = useCallback(async () => {
        if (!rawDescription.trim()) {
            return { success: false, reason: "empty" as const };
        }

        setIsAnalyzing(true);
        try {
            const result = await parseJobDescription(rawDescription);

            // Auto-fill fields
            if (result.title) form.setValue("title", result.title);
            if (result.department) form.setValue("department", result.department);
            if (result.description) form.setValue("description", result.description);
            if (result.type) form.setValue("type", result.type);

            return { success: true as const };
        } catch (error) {
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
