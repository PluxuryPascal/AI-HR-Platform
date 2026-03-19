"use client"

import { useTranslations } from "next-intl"
import { Trash2, Building2 } from "lucide-react"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { useGetDepartments, useDeleteDepartment } from "../api/use-departments"
import { EmptyState } from "@/components/shared/empty-state"
import { CreateDepartmentDialog } from "./create-department-dialog"

export function DepartmentsList() {
    const t = useTranslations("Departments")
    const { data: departments = [], isLoading } = useGetDepartments()
    const { mutate: deleteDepartment } = useDeleteDepartment()

    if (isLoading) {
        return <div className="p-8 text-center">Загрузка...</div>
    }

    if (departments.length === 0) {
        return (
            <EmptyState
                icon={Building2}
                title={t("empty.title")}
                description={t("empty.desc")}
                action={<CreateDepartmentDialog />}
            />
        )
    }

    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>{t("table.name")}</TableHead>
                        <TableHead className="w-[100px] text-right">{t("table.actions")}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {departments.map((dept) => (
                        <TableRow key={dept.id}>
                            <TableCell className="font-medium">{dept.name}</TableCell>
                            <TableCell className="text-right">
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => deleteDepartment(dept.id)}
                                    className="text-red-500 hover:text-red-700 hover:bg-red-100 dark:hover:bg-red-900/30"
                                    title={t("table.delete")}
                                >
                                    <Trash2 className="h-4 w-4" />
                                    <span className="sr-only">{t("table.delete")}</span>
                                </Button>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}
