"use client";

import { useTranslations } from "next-intl";
import { useEffect } from "react";
import { AlertTriangle, RefreshCcw } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function Error({
    error,
    reset,
}: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    const t = useTranslations("Errors.Global");

    useEffect(() => {
        console.error(error);
    }, [error]);

    return (
        <div className="flex h-[calc(100vh-4rem)] flex-col items-center justify-center gap-4 text-center">
            <div className="rounded-full bg-destructive/10 p-4">
                <AlertTriangle className="h-10 w-10 text-destructive" />
            </div>
            <div className="space-y-2">
                <h2 className="text-2xl font-bold tracking-tight">{t("title")}</h2>
                <p className="text-muted-foreground">{t("desc")}</p>
            </div>
            <Button onClick={reset} variant="outline" className="gap-2">
                <RefreshCcw className="h-4 w-4" />
                {t("retry")}
            </Button>
        </div>
    );
}
