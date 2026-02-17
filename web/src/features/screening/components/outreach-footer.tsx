"use client";

import { useTranslations } from "next-intl";
import { Copy, Mail, RefreshCcw } from "lucide-react";
import { SheetFooter } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";

interface OutreachFooterProps {
    content: string;
    isGenerating: boolean;
    isBulk: boolean;
    onRegenerate: () => void;
    onSend: () => void;
}

export function OutreachFooter({
    content,
    isGenerating,
    isBulk,
    onRegenerate,
    onSend,
}: OutreachFooterProps) {
    const t = useTranslations("Outreach");

    return (
        <SheetFooter className="flex-col sm:flex-row gap-2 mt-auto">
            <Button variant="outline" onClick={onRegenerate} disabled={isGenerating}>
                <RefreshCcw className="w-4 h-4 mr-2" />
                {t("regenerate")}
            </Button>
            <div className="flex-1" />
            <Button variant="secondary" onClick={() => navigator.clipboard.writeText(content)}>
                <Copy className="w-4 h-4 mr-2" />
                {t("copy")}
            </Button>
            <Button onClick={onSend}>
                <Mail className="w-4 h-4 mr-2" />
                {isBulk ? t("bulk.sendAll") : t("sendBtn")}
            </Button>
        </SheetFooter>
    );
}
