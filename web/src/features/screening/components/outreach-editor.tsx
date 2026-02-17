"use client";

import { useTranslations } from "next-intl";
import { Sparkles } from "lucide-react";
import {
    Tabs,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";

interface OutreachEditorProps {
    tone: string;
    onToneChange: (tone: string) => void;
    content: string;
    onContentChange: (content: string) => void;
    isGenerating: boolean;
    isRejection: boolean;
}

export function OutreachEditor({
    tone,
    onToneChange,
    content,
    onContentChange,
    isGenerating,
    isRejection,
}: OutreachEditorProps) {
    const t = useTranslations("Outreach");

    return (
        <div className="space-y-4 h-full flex flex-col">
            <div className="flex items-center justify-between">
                <Label>{t("tone.label")}</Label>
                <Tabs value={tone} onValueChange={onToneChange} className="w-auto">
                    <TabsList>
                        <TabsTrigger value="professional">{t("tone.professional")}</TabsTrigger>
                        <TabsTrigger value="friendly">{t("tone.friendly")}</TabsTrigger>
                        <TabsTrigger value="brief">{t("tone.brief")}</TabsTrigger>
                    </TabsList>
                </Tabs>
            </div>

            <div className="relative flex-1">
                <Textarea
                    value={content}
                    onChange={(e) => onContentChange(e.target.value)}
                    className="h-full min-h-[300px] resize-none p-4 font-mono text-sm leading-relaxed"
                    placeholder={t("placeholder")}
                />
                {isGenerating && (
                    <div className="absolute inset-0 bg-background/50 backdrop-blur-[1px] flex items-center justify-center">
                        <div className="flex items-center gap-2 text-primary animate-pulse">
                            <Sparkles className="w-5 h-5" />
                            <span>{t("generating")}</span>
                        </div>
                    </div>
                )}
            </div>

            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Sparkles className="w-3 h-3" />
                <span>{t("reasoning.prefix")} </span>
                <span className="font-medium text-foreground">
                    {isRejection ? t("mockReason.rejection") : t("mockReason.invitation")}
                </span>
            </div>
        </div>
    );
}
