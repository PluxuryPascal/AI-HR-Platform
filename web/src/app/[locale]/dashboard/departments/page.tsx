import { getTranslations } from "next-intl/server"
import { DepartmentsList } from "@/features/departments/components/departments-list"
import { CreateDepartmentDialog } from "@/features/departments/components/create-department-dialog"

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) {
    const { locale } = await params
    const t = await getTranslations({ locale, namespace: "Departments" })
    return {
        title: t("pageTitle"),
    }
}

export default async function DepartmentsPage() {
    const t = await getTranslations("Departments")

    return (
        <div className="flex flex-col gap-6 w-full max-w-5xl mx-auto h-full p-4 md:p-6 lg:p-8 pt-8 md:pt-10">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight mb-2">
                        {t("title")}
                    </h1>
                    <p className="text-muted-foreground text-lg">
                        {t("description")}
                    </p>
                </div>
                <CreateDepartmentDialog />
            </div>

            <DepartmentsList />
        </div>
    )
}
