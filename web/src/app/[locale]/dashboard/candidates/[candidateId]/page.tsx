"use client";

import {
    ResizableHandle,
    ResizablePanel,
    ResizablePanelGroup,
} from "@/components/ui/resizable";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ArrowLeft, Loader2 } from "lucide-react";
import { Link } from "@/i18n/routing"; // Fixed Link import
import { useTranslations } from "next-intl";
import { useGetCandidate } from "@/features/candidates/api/use-get-candidate";
import { use, useState } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import dynamic from "next/dynamic";

const PDFViewer = dynamic(
    () => import("@/features/candidates/components/pdf-viewer").then((mod) => mod.PDFViewer),
    { 
        ssr: false,
        loading: () => (
            <div className="flex-1 flex items-center justify-center bg-background/50 backdrop-blur-sm">
                <div className="flex flex-col items-center gap-3">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                </div>
            </div>
        )
    }
);
import { AIAnalysisTab } from "@/features/candidates/components/ai-analysis-tab";
import { AIChatTab } from "@/features/candidates/components/ai-chat-tab";
import { useParams } from "next/navigation";
import { WidgetErrorBoundary } from "@/components/shared/widget-error-boundary";
import { InterviewGuide } from "@/features/candidates/components/interview-guide";
import { MatchScoreBadge } from "@/features/candidates/components/match-score-badge";


import { OutreachDrawer } from "@/features/screening/components/outreach-drawer";


export default function CandidatePage({ params }: { params: Promise<{ candidateId: string; locale: string }> }) {
    const { candidateId } = use(params);
    const t = useTranslations("CandidateProfile");
    const [isOutreachOpen, setIsOutreachOpen] = useState(false);

    const { data: detail, isLoading, error } = useGetCandidate(candidateId);

    if (isLoading) {
        return (
            <div className="flex-1 space-y-4 p-8 pt-6">
                <Skeleton className="h-20 w-full" />
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                    <Skeleton className="col-span-4 h-[600px]" />
                    <Skeleton className="col-span-3 h-[600px]" />
                </div>
            </div>
        );
    }

    if (error || !detail) {
        return (
            <div className="flex-1 flex items-center justify-center p-8">
                <div className="text-center space-y-2">
                    <h2 className="text-2xl font-bold tracking-tight">Кандидат не найден</h2>
                    <p className="text-muted-foreground">Не удалось загрузить данные кандидата.</p>
                </div>
            </div>
        );
    }

    const { candidate, profile, score, factors } = detail;

    // Map backend data to UI format
    const uiProfile = {
        personal: {
            name: `${candidate.first_name || ""} ${candidate.last_name || ""}`.trim() || "Unknown Candidate",
            role: "Software Engineer", // This should come from Job, but for now hardcode or use parsed data
            location: candidate.location || "Remote",
            email: candidate.email || "",
            phone: "+1 (555) 000-0000", // Mock if not in DB
        },
        aiAnalysis: {
            score: score?.match_score || 0,
            summary: candidate.parsed_text?.slice(0, 500) + "..." || "",
            scoreBreakdown: factors?.map(f => ({
                id: f.id,
                text: f.description,
                impact: f.impact > 0 ? `+${f.impact}` : `${f.impact}`,
                type: f.type
            })) || [],
            strengths: factors?.filter(f => f.type === 'positive').map(f => f.description) || [],
            weaknesses: factors?.filter(f => f.type === 'negative').map(f => f.description) || [],
            skills: candidate.skills || [],
        },
        pdfUrl: candidate.resume_url || "",
    };

    return (
        <div className="h-full flex flex-col bg-background">
            {/* Header */}
            <div className="flex items-center justify-between px-6 py-3 border-b bg-card shrink-0">
                <div className="flex items-center gap-4">
                    <Button variant="ghost" size="icon" asChild>
                        <Link href="/dashboard/candidates">
                            <ArrowLeft className="w-4 h-4" />
                        </Link>
                    </Button>
                    <div>
                        <div className="flex items-center gap-2">
                            <h1 className="font-semibold text-lg">{uiProfile.personal.name}</h1>
                            <Badge variant="outline" className="bg-primary/5 text-primary border-primary/20">
                                {uiProfile.personal.role}
                            </Badge>
                            <MatchScoreBadge
                                score={uiProfile.aiAnalysis.score}
                                breakdown={uiProfile.aiAnalysis.scoreBreakdown}
                                className="ml-2"
                            />
                        </div>
                        <p className="text-xs text-muted-foreground">
                            Applied 2 days ago • {uiProfile.personal.location}
                        </p>
                    </div>
                </div>
                <div className="flex items-center gap-2">
                </div>
            </div>

            {/* Main Content - Split View */}
            <div className="flex-1 overflow-hidden bg-background">
                <ResizablePanelGroup direction="horizontal" className="h-full w-full rounded-none border-0">

                    <ResizablePanel defaultSize={50} minSize={30}>
                        <div className="h-full bg-background flex flex-col transition-colors duration-300">
                            {/* PDF Header - Synchronized with Right Panel */}
                            <div className="h-12 bg-background/80 backdrop-blur-md border-b border-border/50 px-4 flex items-center justify-between shrink-0 z-10">
                                <span className="font-medium text-sm text-foreground/80 pt-2">{t("tabs.resume")}</span>
                                <div className="flex items-center gap-2 pt-2">
                                    <Badge variant="outline" className="text-[10px] h-5 bg-transparent border-foreground/20 text-muted-foreground">PDF</Badge>
                                </div>
                            </div>
                            <div className="flex-1 overflow-hidden relative">
                                <WidgetErrorBoundary>
                                    <PDFViewer url={uiProfile.pdfUrl} />
                                </WidgetErrorBoundary>
                            </div>
                        </div>
                    </ResizablePanel>

                    <ResizableHandle withHandle />

                    {/* Right Panel - AI Analysis & Chat */}
                    <ResizablePanel defaultSize={50} minSize={30}>
                        <div className="h-full flex flex-col bg-background">
                            <Tabs defaultValue="analysis" className="flex-1 flex flex-col h-full">
                                <div className="h-12 border-b border-border/50 px-4 bg-background/80 backdrop-blur-md shrink-0 flex items-center">
                                    <TabsList className="h-full w-full justify-start bg-transparent p-0 gap-6">
                                        <TabsTrigger
                                            value="analysis"
                                            className="h-full rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-2 pt-2"
                                        >
                                            {t("tabs.overview")}
                                        </TabsTrigger>
                                        <TabsTrigger
                                            value="chat"
                                            className="h-full rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-2 pt-2"
                                        >
                                            {t("tabs.chat")}
                                        </TabsTrigger>
                                        <TabsTrigger
                                            value="interview"
                                            className="h-full rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-2 pt-2"
                                        >
                                            {t("tabs.interview")}
                                        </TabsTrigger>
                                    </TabsList>
                                </div>

                                <div className="flex-1 overflow-hidden">
                                    <TabsContent value="analysis" className="h-full m-0 p-0 border-none select-text data-[state=active]:flex data-[state=active]:flex-col overflow-y-auto">
                                        <AIAnalysisTab data={uiProfile.aiAnalysis} onOutreach={() => setIsOutreachOpen(true)} />
                                    </TabsContent>
                                    <TabsContent value="chat" className="h-full m-0 p-0 border-none data-[state=active]:flex data-[state=active]:flex-col">
                                        <WidgetErrorBoundary>
                                            <AIChatTab />
                                        </WidgetErrorBoundary>
                                    </TabsContent>
                                    <TabsContent value="interview" className="h-full m-0 p-0 border-none select-text data-[state=active]:flex data-[state=active]:flex-col overflow-y-auto">
                                        <InterviewGuide
                                            matchScore={uiProfile.aiAnalysis.score}
                                            candidateId={candidateId}
                                        />
                                    </TabsContent>
                                </div>
                            </Tabs>
                        </div>
                    </ResizablePanel>

                </ResizablePanelGroup>
            </div>

            <OutreachDrawer
                isOpen={isOutreachOpen}
                onClose={() => setIsOutreachOpen(false)}
                candidate={{
                    id: candidate.id,
                    name: uiProfile.personal.name,
                    role: uiProfile.personal.role,
                    score: uiProfile.aiAnalysis.score,
                    email: uiProfile.personal.email,
                    matchSummary: uiProfile.aiAnalysis.summary,
                    scoreBreakdown: uiProfile.aiAnalysis.scoreBreakdown as any,
                }}
                type="invitation"
            />
        </div>
    );
}
