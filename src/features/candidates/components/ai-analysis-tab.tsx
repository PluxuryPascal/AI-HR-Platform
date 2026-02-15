import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { useTranslations } from "next-intl";
import { CheckCircle2, XCircle, Zap } from "lucide-react";
import { Profile } from "../utils/mock-profile";

import { Button } from "@/components/ui/button";
import { MessageSquare } from "lucide-react";

interface AIAnalysisTabProps {
    data: Profile["aiAnalysis"];
    onOutreach?: () => void;
}

export const AIAnalysisTab = ({ data, onOutreach }: AIAnalysisTabProps) => {
    const t = useTranslations("CandidateProfile.analysis");
    const tOutreach = useTranslations("Outreach");

    return (
        <ScrollArea className="h-full">
            <div className="p-6 space-y-8">
                {/* Score Section */}
                <div className="flex flex-col items-center justify-center space-y-4">
                    <div className="relative flex items-center justify-center w-32 h-32">
                        <svg className="w-full h-full transform -rotate-90">
                            <circle
                                cx="64"
                                cy="64"
                                r="60"
                                stroke="currentColor"
                                strokeWidth="8"
                                fill="transparent"
                                className="text-gray-100"
                            />
                            <circle
                                cx="64"
                                cy="64"
                                r="60"
                                stroke="currentColor"
                                strokeWidth="8"
                                fill="transparent"
                                strokeDasharray={377}
                                strokeDashoffset={377 - (377 * data.score) / 100}
                                className={`text-blue-600 transition-all duration-1000 ease-out`}
                            />
                        </svg>
                        <span className="absolute text-3xl font-bold">{data.score}</span>
                    </div>
                    <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
                        {t("matchScore")}
                    </p>
                </div>

                <Separator />

                {/* Summary */}
                <div className="space-y-3">
                    <h3 className="text-lg font-semibold flex items-center">
                        <Zap className="w-5 h-5 mr-2 text-yellow-500" />
                        {t("summary")}
                    </h3>
                    <p className="text-muted-foreground leading-relaxed">
                        {data.summary}
                    </p>
                    <Button
                        onClick={onOutreach}
                        className="w-full mt-4"
                        variant="outline"
                    >
                        <MessageSquare className="w-4 h-4 mr-2" />
                        {tOutreach("action.sendResponse")}
                    </Button>
                </div>

                {/* Strengths */}
                <div className="space-y-3">
                    <h3 className="text-lg font-semibold flex items-center text-green-700 dark:text-green-500">
                        <CheckCircle2 className="w-5 h-5 mr-2" />
                        {t("strengths")}
                    </h3>
                    <ul className="space-y-2">
                        {data.strengths.map((strength, i) => (
                            <li key={i} className="flex items-start text-sm">
                                <span className="w-1.5 h-1.5 rounded-full bg-green-500 mt-2 mr-2 shrink-0" />
                                {strength}
                            </li>
                        ))}
                    </ul>
                </div>

                {/* Weaknesses */}
                <div className="space-y-3">
                    <h3 className="text-lg font-semibold flex items-center text-red-700 dark:text-red-500">
                        <XCircle className="w-5 h-5 mr-2" />
                        {t("weaknesses")}
                    </h3>
                    <ul className="space-y-2">
                        {data.weaknesses.map((weakness, i) => (
                            <li key={i} className="flex items-start text-sm">
                                <span className="w-1.5 h-1.5 rounded-full bg-red-500 mt-2 mr-2 shrink-0" />
                                {weakness}
                            </li>
                        ))}
                    </ul>
                </div>

                <Separator />

                {/* Skills */}
                <div className="space-y-3">
                    <h3 className="text-lg font-semibold">{t("skills")}</h3>
                    <div className="flex flex-wrap gap-2">
                        {data.skills.map((skill) => (
                            <Badge key={skill} variant="secondary">
                                {skill}
                            </Badge>
                        ))}
                    </div>
                </div>
            </div>
        </ScrollArea>
    );
};
