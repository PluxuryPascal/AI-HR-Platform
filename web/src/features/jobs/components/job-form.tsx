"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslations } from "next-intl"
import { toast } from "sonner"
import { Loader2, Briefcase } from "lucide-react"
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
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"

import {
    jobSchema,
    JobFormValues,
    Department,
    JobType,
} from "@/features/jobs/schemas/job-schema"
import { SmartPasteCard } from "./smart-paste-card"

export function JobForm() {
    const t = useTranslations("JobWizard")
    const router = useRouter()

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

    const onSubmit = async (data: JobFormValues) => {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 1000))
        console.log("Job Created:", data)
        toast.success(t("successMessage"))
    }

    return (
        <div className="space-y-6">
            {/* Smart Paste Section */}
            <SmartPasteCard form={form} />

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
