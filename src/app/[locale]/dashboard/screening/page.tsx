import { KanbanBoard } from "@/features/screening/components/kanban-board";
import { getTranslations } from "next-intl/server";

export default async function ScreeningPage() {
    const t = await getTranslations("Screening");

    return (
        <div className="flex flex-col h-[calc(100vh-4rem)]">
            <div className="flex items-center justify-between px-6 py-4 border-b">
                <h1 className="text-2xl font-bold tracking-tight">{t("title")}</h1>
            </div>
            <div className="flex-1 overflow-hidden p-6">
                <KanbanBoard />
            </div>
        </div>
    );
}
