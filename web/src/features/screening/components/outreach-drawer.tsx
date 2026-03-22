"use client";

import { useTranslations } from "next-intl";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { CandidateCard } from "../types";
import { OutreachEditor } from "./outreach-editor";
import { OutreachFooter } from "./outreach-footer";
import { useResizablePanel } from "../hooks/use-resizable-panel";
import { useOutreachGenerator } from "../hooks/use-outreach-generator";
import { useSendEmail } from "@/features/candidates/api/use-outreach";

interface OutreachDrawerProps {
    isOpen: boolean;
    onClose: () => void;
    candidate: CandidateCard | CandidateCard[] | null;
    type: "rejection" | "invitation";
    onSend?: (content: string) => void;
}

export function OutreachDrawer({
    isOpen,
    onClose,
    candidate,
    type,
    onSend,
}: OutreachDrawerProps) {
    const t = useTranslations("Outreach");
    const { width, startResizing } = useResizablePanel(600);
    const {
        tone,
        setTone,
        content,
        setContent,
        communicationId,
        isGenerating,
        isBulk,
        candidates,
        handleGenerate,
    } = useOutreachGenerator({ isOpen, candidate, type });

    const sendEmail = useSendEmail();

    const handleSend = async () => {
        if (!communicationId || !content) return;
        
        // Subject is usually the first line or we could parse it more robustly.
        // For now, let's assume content starts with subject if non-bulk.
        const lines = content.split("\n");
        const subject = lines[0];
        const body = lines.slice(1).join("\n").trim();

        try {
            await sendEmail.mutateAsync({
                communicationId,
                subject,
                body,
            });
            if (onSend) onSend(content);
            onClose();
        } catch (error) {
            // Error handled by hook toast
        }
    };

    if (candidates.length === 0) return null;

    const isRejection = type === "rejection";

    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent
                className="w-full sm:max-w-none p-0 flex flex-row overflow-hidden"
                style={{ width: `${width}px` }}
            >
                {/* Resize Handle */}
                <div
                    className="w-1.5 h-full bg-border hover:bg-primary/20 cursor-ew-resize transition-colors flex flex-col justify-center items-center group touch-none"
                    onMouseDown={startResizing}
                >
                    <div className="h-8 w-1 rounded-full bg-muted-foreground/30 group-hover:bg-primary/50" />
                </div>

                <div className="flex-1 flex flex-col h-full overflow-y-auto p-6">
                    <SheetHeader className="space-y-4">
                        <div className="flex items-center gap-2">
                            <Badge variant={isRejection ? "destructive" : "default"}>
                                {isRejection ? t("badges.rejection") : t("badges.invitation")}
                            </Badge>
                            <span className="text-sm text-muted-foreground">
                                {!isBulk && <>{t("to")}: {candidates[0].name}</>}
                                {isBulk && `${candidates.length} candidates selected`}
                            </span>
                        </div>
                        <SheetTitle>
                            {isBulk
                                ? t("bulk.title")
                                : (isRejection ? t("title.rejected") : t("title.invited"))}
                        </SheetTitle>
                        <SheetDescription>
                            {isBulk ? t("bulk.description") : t("description")}
                        </SheetDescription>
                    </SheetHeader>

                    <div className="py-6 space-y-6 flex-1">
                        <OutreachEditor
                            tone={tone}
                            onToneChange={setTone}
                            content={content}
                            onContentChange={setContent}
                            isGenerating={isGenerating}
                            isRejection={isRejection}
                        />
                    </div>

                    <OutreachFooter
                        content={content}
                        isGenerating={isGenerating}
                        isBulk={isBulk}
                        onRegenerate={handleGenerate}
                        onSend={handleSend}
                    />
                </div>
            </SheetContent>
        </Sheet>
    );
}
