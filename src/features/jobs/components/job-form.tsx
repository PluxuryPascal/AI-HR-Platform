"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslations } from "next-intl"
import { toast } from "sonner"
import { Loader2, Sparkles, Briefcase, Building, FileText, CheckCircle2 } from "lucide-react"
import { useRouter } from "next/navigation"
import { cn } from "@/lib/utils"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Label } from "@/components/ui/label"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"

import {
    jobSchema,
    JobFormValues,
    Department,
    JobType,
} from "@/features/jobs/schemas/job-schema"
import { parseJobDescription } from "@/features/jobs/utils/mock-ai-parser"

export function JobForm() {
    const t = useTranslations("JobWizard")
    const router = useRouter()
    const [isAnalyzing, setIsAnalyzing] = useState(false)
    const [rawDescription, setRawDescription] = useState("")

    const form = useForm<JobFormValues>({
        resolver: zodResolver(jobSchema),
        defaultValues: {
            title: "",
            department: undefined,
            description: "",
            requirements: [],
            type: JobType.Onsite,
        },
    })

    const onAnalyze = async () => {
        if (!rawDescription.trim()) {
            toast.error(t("smartPastePlaceholder"))
            return
        }

        setIsAnalyzing(true)
        try {
            const result = await parseJobDescription(rawDescription)

            // Auto-fill fields
            if (result.title) form.setValue("title", result.title)
            if (result.department) form.setValue("department", result.department)
            if (result.description) form.setValue("description", result.description)
            if (result.type) form.setValue("type", result.type)

            toast.success(t("aiSuccess"), {
                description: t("aiSuccessDesc"),
            })
        } catch (error) {
            toast.error(t("aiError"), {
                description: t("aiErrorDesc"),
            })
        } finally {
            setIsAnalyzing(false)
        }
    }

    const onSubmit = async (data: JobFormValues) => {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 1000))
        console.log("Job Created:", data)
        toast.success(t("successMessage"))
    }

    return (
        <div className="space-y-6">
            {/* Smart Paste Section */}
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
                        onClick={onAnalyze}
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

            <Separator className="my-6" />

            {/* Main Form */}
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                <Card>
                    <CardContent className="p-8 space-y-8">

                        {/* Title & Department Row */}
                        <div className="grid gap-6 md:grid-cols-2">
                            <div className="space-y-3">
                                <Label htmlFor="title" className="mb-2 block">{t("fields.title")}</Label>
                                <div className="relative">
                                    <Briefcase className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="title"
                                        placeholder={t("titlePlaceholder")}
                                        className="pl-9"
                                        {...form.register("title")}
                                    />
                                </div>
                                {form.formState.errors.title && (
                                    <p className="text-sm text-red-500">{form.formState.errors.title.message}</p>
                                )}
                            </div>

                            <div className="space-y-3">
                                <Label htmlFor="department" className="mb-2 block">{t("fields.department")}</Label>
                                <Select
                                    onValueChange={(val) => form.setValue("department", val as any)}
                                    defaultValue={form.getValues("department")}
                                    value={form.watch("department")}
                                >
                                    <SelectTrigger className="w-full">
                                        <SelectValue placeholder={t("departmentPlaceholder")} />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {Object.values(Department).map((dept) => (
                                            <SelectItem key={dept} value={dept}>
                                                {t(`departments.${dept}`)}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                                {form.formState.errors.department && (
                                    <p className="text-sm text-red-500">{form.formState.errors.department.message}</p>
                                )}
                            </div>
                        </div>

                        {/* Job Type */}
                        <div className="space-y-3">
                            <Label className="mb-2 block">{t("fields.type")}</Label>
                            <div className="flex gap-4">
                                {Object.values(JobType).map((type) => (
                                    <label
                                        key={type}
                                        className={cn(
                                            "flex cursor-pointer items-center gap-2 rounded-lg border p-3 transition-all hover:bg-accent",
                                            form.watch("type") === type ? "border-primary bg-primary/5 ring-1 ring-primary" : "border-input"
                                        )}
                                    >
                                        <input
                                            type="radio"
                                            value={type}
                                            className="sr-only"
                                            {...form.register("type")}
                                        />
                                        <div className={cn(
                                            "h-4 w-4 rounded-full border border-primary flex items-center justify-center",
                                            form.watch("type") === type ? "bg-primary" : "bg-transparent"
                                        )}>
                                            {form.watch("type") === type && <div className="h-2 w-2 rounded-full bg-white" />}
                                        </div>
                                        <span className="text-sm font-medium">{t(`jobTypes.${type}`)}</span>
                                    </label>
                                ))}
                            </div>
                        </div>

                        {/* Description */}
                        <div className="space-y-3">
                            <Label htmlFor="description" className="mb-2 block">{t("fields.description")}</Label>
                            <Textarea
                                id="description"
                                placeholder={t("descriptionPlaceholder")}
                                className="min-h-[200px]"
                                {...form.register("description")}
                            />
                            {form.formState.errors.description && (
                                <p className="text-sm text-red-500">{form.formState.errors.description.message}</p>
                            )}
                        </div>

                    </CardContent>
                </Card>

                <div className="flex justify-end gap-4">
                    <Button type="button" variant="outline" onClick={() => router.back()}>
                        {t("cancel")}
                    </Button>
                    <Button type="submit" disabled={form.formState.isSubmitting} className="min-w-[140px]">
                        {form.formState.isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        {t("submitBtn")}
                    </Button>
                </div>
            </form>
        </div>
    )
}
