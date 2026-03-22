"use client"

import { useForm, useFieldArray } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslations } from "next-intl"
import { toast } from "sonner"
import { Loader2, Briefcase, Plus, Trash2, Banknote } from "lucide-react"
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
    JobType,
} from "@/features/jobs/schemas/job-schema"
import { useGetDepartments } from "@/features/departments/api/use-departments"
import { useCreateJob, CreateJobPayload } from "@/features/jobs/api/use-create-job"
import { SmartPasteCard } from "./smart-paste-card"

export function JobForm() {
    const t = useTranslations("JobWizard")
    const router = useRouter()
    const { data: departments = [], isLoading: isLoadingDepts } = useGetDepartments()
    const createJobMutation = useCreateJob()

    const form = useForm<JobFormValues>({
        resolver: zodResolver(jobSchema),
        defaultValues: {
            title: "",
            department: "",
            description: "",
            requirements: [""],
            type: JobType.Onsite,
            salary_min: 0,
            salary_max: 0,
            currency: "RUB",
        },
    })

    const { fields, append, remove } = useFieldArray({
        control: form.control,
        name: "requirements",
    })

    const onSubmit = async (data: JobFormValues) => {
        const payload: CreateJobPayload = {
            title: data.title,
            department_id: data.department,
            description: data.description,
            requirements: data.requirements.filter(r => r.trim() !== ""),
            work_format: (data.type === JobType.Remote ? "remote" : data.type === JobType.Onsite ? "office" : "hybrid") as "remote" | "office" | "hybrid",
            salary_min: data.salary_min,
            salary_max: data.salary_max,
            currency: data.currency,
        }

        try {
            await createJobMutation.mutateAsync(payload)
            router.push("/dashboard/jobs")
        } catch (error) {
            // Error is handled in the mutation
        }
    }

    if (isLoadingDepts) {
        return (
            <div className="flex items-center justify-center p-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        )
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
                                        disabled={createJobMutation.isPending}
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
                                    onValueChange={(val) => form.setValue("department", val)}
                                    defaultValue={form.getValues("department")}
                                    value={form.watch("department")}
                                    disabled={createJobMutation.isPending}
                                >
                                    <SelectTrigger className="w-full">
                                        <SelectValue placeholder={t("departmentPlaceholder")} />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {departments.map((dept) => (
                                            <SelectItem key={dept.id} value={dept.id}>
                                                {dept.name}
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
                                            form.watch("type") === type ? "border-primary bg-primary/5 ring-1 ring-primary" : "border-input",
                                            createJobMutation.isPending && "opacity-50 cursor-not-allowed"
                                        )}
                                    >
                                        <input
                                            type="radio"
                                            value={type}
                                            className="sr-only"
                                            disabled={createJobMutation.isPending}
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

                        {/* Salary & Currency */}
                        <div className="grid gap-6 md:grid-cols-3">
                            <div className="space-y-3">
                                <Label htmlFor="salary_min" className="mb-2 block">{t("fields.salaryMin")}</Label>
                                <div className="relative">
                                    <Banknote className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="salary_min"
                                        type="number"
                                        className="pl-9"
                                        disabled={createJobMutation.isPending}
                                        {...form.register("salary_min", { valueAsNumber: true })}
                                    />
                                </div>
                            </div>
                            <div className="space-y-3">
                                <Label htmlFor="salary_max" className="mb-2 block">{t("fields.salaryMax")}</Label>
                                <div className="relative">
                                    <Banknote className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="salary_max"
                                        type="number"
                                        className="pl-9"
                                        disabled={createJobMutation.isPending}
                                        {...form.register("salary_max", { valueAsNumber: true })}
                                    />
                                </div>
                            </div>
                            <div className="space-y-3">
                                <Label htmlFor="currency" className="mb-2 block">{t("fields.currency")}</Label>
                                <Input
                                    id="currency"
                                    placeholder="RUB, USD, EUR..."
                                    disabled={createJobMutation.isPending}
                                    {...form.register("currency")}
                                />
                            </div>
                        </div>

                        {/* Description */}
                        <div className="space-y-3">
                            <Label htmlFor="description" className="mb-2 block">{t("fields.description")}</Label>
                            <Textarea
                                id="description"
                                placeholder={t("descriptionPlaceholder")}
                                className="min-h-[200px]"
                                disabled={createJobMutation.isPending}
                                {...form.register("description")}
                            />
                            {form.formState.errors.description && (
                                <p className="text-sm text-red-500">{form.formState.errors.description.message}</p>
                            )}
                        </div>

                        {/* Requirements */}
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <Label className="block">{t("fields.requirements")}</Label>
                                <Button
                                    type="button"
                                    variant="outline"
                                    size="sm"
                                    onClick={() => append("")}
                                    disabled={createJobMutation.isPending}
                                >
                                    <Plus className="h-4 w-4 mr-2" />
                                    Добавить
                                </Button>
                            </div>
                            <div className="space-y-3">
                                {fields.map((field, index) => (
                                    <div key={field.id} className="flex gap-2">
                                        <Input
                                            {...form.register(`requirements.${index}`)}
                                            placeholder={`Требование #${index + 1}`}
                                            disabled={createJobMutation.isPending}
                                        />
                                        <Button
                                            type="button"
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => remove(index)}
                                            disabled={createJobMutation.isPending || fields.length === 1}
                                            className="text-red-500 hover:text-red-600 hover:bg-red-50"
                                        >
                                            <Trash2 className="h-4 w-4" />
                                        </Button>
                                    </div>
                                ))}
                            </div>
                            {form.formState.errors.requirements && (
                                <p className="text-sm text-red-500">{form.formState.errors.requirements.message}</p>
                            )}
                        </div>

                    </CardContent>
                </Card>

                <div className="flex justify-end gap-4">
                    <Button type="button" variant="outline" onClick={() => router.back()} disabled={createJobMutation.isPending}>
                        {t("cancel")}
                    </Button>
                    <Button type="submit" disabled={createJobMutation.isPending} className="min-w-[140px]">
                        {createJobMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        {t("submitBtn")}
                    </Button>
                </div>
            </form>
        </div>
    )
}
