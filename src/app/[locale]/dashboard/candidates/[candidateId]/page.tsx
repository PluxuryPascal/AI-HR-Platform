"use client";

import {
    ResizableHandle,
    ResizablePanel,
    ResizablePanelGroup,
} from "@/components/ui/resizable";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ArrowLeft, Share2, MoreHorizontal } from "lucide-react";
import { Link } from "@/i18n/routing"; // Fixed Link import
import { useTranslations } from "next-intl";
import { getCandidateProfile } from "@/features/candidates/utils/mock-profile";
import { PDFViewer } from "@/features/candidates/components/pdf-viewer";
import { AIAnalysisTab } from "@/features/candidates/components/ai-analysis-tab";
import { AIChatTab } from "@/features/candidates/components/ai-chat-tab";
import { useParams } from "next/navigation";
import { WidgetErrorBoundary } from "@/components/shared/widget-error-boundary";
import { InterviewGuide } from "@/features/candidates/components/interview-guide";

export default function CandidatePage() {
    const params = useParams();
    const t = useTranslations("CandidateProfile");
    const profile = getCandidateProfile(params.candidateId as string);

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
                            <h1 className="font-semibold text-lg">{profile.personal.name}</h1>
                            <Badge variant="outline" className="text-xs font-normal">
                                {profile.personal.role}
                            </Badge>
                        </div>
                        <p className="text-xs text-muted-foreground">
                            Applied 2 days ago • {profile.personal.location}
                        </p>
                    </div>
                </div>
                <div className="flex items-center gap-2">
                    <Button variant="outline" size="sm">
                        <Share2 className="w-4 h-4 mr-2" />
                        Share
                    </Button>
                    <Button variant="ghost" size="icon">
                        <MoreHorizontal className="w-4 h-4" />
                    </Button>
                </div>
            </div>

            {/* Main Content - Split View */}
            <div className="flex-1 overflow-hidden">
                <ResizablePanelGroup direction="horizontal" className="h-full w-full rounded-none border-0">

                    {/* Left Panel - PDF Viewer */}
                    <ResizablePanel defaultSize={50} minSize={30}>
                        <div className="h-full bg-gray-100 flex flex-col">
                            <div className="bg-white border-b px-4 py-2 flex items-center justify-between text-sm text-gray-500 shadow-sm z-10 shrink-0">
                                <span className="font-medium text-gray-700">{t("tabs.resume")}</span>
                                <div className="flex items-center gap-2">
                                    <Badge variant="secondary" className="text-[10px] h-5">PDF</Badge>
                                </div>
                            </div>
                            <div className="flex-1 overflow-hidden relative">
                                <WidgetErrorBoundary>
                                    <PDFViewer url={profile.pdfUrl} />
                                </WidgetErrorBoundary>
                            </div>
                        </div>
                    </ResizablePanel>

                    <ResizableHandle withHandle />

                    {/* Right Panel - AI Analysis & Chat */}
                    <ResizablePanel defaultSize={50} minSize={30}>
                        <div className="h-full flex flex-col bg-background">
                            <Tabs defaultValue="overview" className="flex-1 flex flex-col h-full">
                                <div className="border-b px-4 bg-background shrink-0">
                                    <TabsList className="h-12 w-full justify-start bg-transparent p-0 gap-6">
                                        <TabsTrigger
                                            value="overview"
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
                                    <TabsContent value="overview" className="h-full m-0 p-0 border-none select-text data-[state=active]:flex data-[state=active]:flex-col">
                                        <AIAnalysisTab data={profile.aiAnalysis} />
                                    </TabsContent>
                                    <TabsContent value="chat" className="h-full m-0 p-0 border-none data-[state=active]:flex data-[state=active]:flex-col">
                                        <WidgetErrorBoundary>
                                            <AIChatTab />
                                        </WidgetErrorBoundary>
                                    </TabsContent>
                                    <TabsContent value="interview" className="h-full m-0 p-0 border-none select-text data-[state=active]:flex data-[state=active]:flex-col overflow-y-auto">
                                        <InterviewGuide matchScore={profile.aiAnalysis.score} />
                                    </TabsContent>
                                </div>
                            </Tabs>
                        </div>
                    </ResizablePanel>

                </ResizablePanelGroup>
            </div>
        </div>
    );
}
