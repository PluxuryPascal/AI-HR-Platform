"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { Copy, Mail, Sparkles, RefreshCcw } from "lucide-react";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetFooter,
} from "@/components/ui/sheet";
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { CandidateCard } from "../types";
import { generateOutreachEmail } from "../utils/outreach-generator";
import { useLocale } from "next-intl";

interface OutreachDrawerProps {
    isOpen: boolean;
    onClose: () => void;
    candidate: CandidateCard | null;
    type: "rejection" | "invitation";
}

export function OutreachDrawer({
    isOpen,
    onClose,
    candidate,
    type,
}: OutreachDrawerProps) {
    const t = useTranslations("Outreach");
    const locale = useLocale();
    const [tone, setTone] = useState("professional");
    const [content, setContent] = useState("");
    const [isGenerating, setIsGenerating] = useState(false);

    useEffect(() => {
        if (isOpen && candidate) {
            handleGenerate();
        }
    }, [isOpen, candidate, type, tone]);

    const handleGenerate = () => {
        if (!candidate) return;
        setIsGenerating(true);

        // Simulate AI delay
        setTimeout(() => {
            const newContent = generateOutreachEmail(candidate, type, tone, locale);
            setContent(newContent);
            setIsGenerating(false);
        }, 600);
    };

    if (!candidate) return null;

    const [width, setWidth] = useState(600);
    const [isResizing, setIsResizing] = useState(false);

    useEffect(() => {
        const handleMouseMove = (e: MouseEvent) => {
            if (!isResizing) return;
            const newWidth = window.innerWidth - e.clientX;
            // Limit width between 400px and 90% of screen width
            if (newWidth > 400 && newWidth < window.innerWidth * 0.9) {
                setWidth(newWidth);
            }
        };

        const handleMouseUp = () => {
            setIsResizing(false);
            document.body.style.cursor = 'default';
        };

        if (isResizing) {
            document.addEventListener('mousemove', handleMouseMove);
            document.addEventListener('mouseup', handleMouseUp);
            document.body.style.cursor = 'ew-resize';
        }

        return () => {
            document.removeEventListener('mousemove', handleMouseMove);
            document.removeEventListener('mouseup', handleMouseUp);
            document.body.style.cursor = 'default';
        };
    }, [isResizing]);

    if (!candidate) return null;

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
                    onMouseDown={(e) => {
                        e.preventDefault();
                        setIsResizing(true);
                    }}
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
                                {t("to")}: {candidate.name}
                            </span>
                        </div>
                        <SheetTitle>
                            {isRejection ? t("title.rejected") : t("title.invited")}
                        </SheetTitle>
                        <SheetDescription>
                            {t("description")}
                        </SheetDescription>
                    </SheetHeader>

                    <div className="py-6 space-y-6 flex-1">
                        <div className="space-y-4 h-full flex flex-col">
                            <div className="flex items-center justify-between">
                                <Label>{t("tone.label")}</Label>
                                <Tabs value={tone} onValueChange={setTone} className="w-auto">
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
                                    onChange={(e) => setContent(e.target.value)}
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
                    </div>

                    <SheetFooter className="flex-col sm:flex-row gap-2 mt-auto">
                        <Button variant="outline" onClick={handleGenerate} disabled={isGenerating}>
                            <RefreshCcw className="w-4 h-4 mr-2" />
                            {t("regenerate")}
                        </Button>
                        <div className="flex-1" />
                        <Button variant="secondary" onClick={() => navigator.clipboard.writeText(content)}>
                            <Copy className="w-4 h-4 mr-2" />
                            {t("copy")}
                        </Button>
                        <Button onClick={onClose}>
                            <Mail className="w-4 h-4 mr-2" />
                            {t("sendBtn")}
                        </Button>
                    </SheetFooter>
                </div>
            </SheetContent>
        </Sheet>
    );
}
