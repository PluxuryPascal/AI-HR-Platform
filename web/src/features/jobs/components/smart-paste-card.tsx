"use client"

import { useTranslations } from "next-intl"
import { toast } from "sonner"
import { Loader2, Sparkles } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { UseFormReturn } from "react-hook-form"
import { JobFormValues } from "@/features/jobs/schemas/job-schema"
import { useJobAiParser } from "@/features/jobs/hooks/use-job-ai-parser"

interface SmartPasteCardProps {
    form: UseFormReturn<JobFormValues>;
}

export function SmartPasteCard({ form }: SmartPasteCardProps) {
    const t = useTranslations("JobWizard")
    const { isAnalyzing, rawDescription, setRawDescription, onAnalyze } = useJobAiParser(form)

    const handleAnalyze = async () => {
        const result = await onAnalyze()
        if (!result.success && result.reason === "empty") {
            toast.error(t("smartPastePlaceholder"))
        } else if (!result.success && result.reason === "error") {
            toast.error(t("aiError"), {
                description: t("aiErrorDesc"),
            })
        } else if (result.success) {
            toast.success(t("aiSuccess"), {
                description: t("aiSuccessDesc"),
            })
        }
    }

    return (
        <Card className="border-indigo-500/20 bg-indigo-50/10 dark:bg-indigo-950/10">
            <CardHeader>
                <CardTitle className="flex items-center gap-2 text-indigo-600 dark:text-indigo-400">
                    <Sparkles className="h-5 w-5" />
                    {t("smartPasteLabel")}
                </CardTitle>
                <CardDescription>
                    {t("smartPastePlaceholder")}
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <Textarea
                    placeholder={t("smartPastePlaceholder")}
                    className="min-h-[120px] resize-y border-indigo-200 focus-visible:ring-indigo-500 dark:border-indigo-800"
                    value={rawDescription}
                    onChange={(e) => setRawDescription(e.target.value)}
                />
                <Button
                    type="button"
                    onClick={handleAnalyze}
                    disabled={isAnalyzing || !rawDescription}
                    className="w-full bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white shadow-lg shadow-indigo-500/20"
                >
                    {isAnalyzing ? (
                        <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            Analyzing...
                        </>
                    ) : (
                        <>
                            <Sparkles className="mr-2 h-4 w-4" />
                            {t("autoFillBtn")}
                        </>
                    )}
                </Button>
            </CardContent>
        </Card>
    )
}
