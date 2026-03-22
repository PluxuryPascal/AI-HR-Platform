import { useEffect, useState } from "react";
import { useLocale } from "next-intl";
import { CandidateCard } from "../types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { useGenerateEmail } from "@/features/candidates/api/use-outreach";

interface UseOutreachGeneratorOptions {
    isOpen: boolean;
    candidate: CandidateCard | CandidateCard[] | null;
    type: "rejection" | "invitation";
}

export function useOutreachGenerator({ isOpen, candidate, type }: UseOutreachGeneratorOptions) {
    const locale = useLocale();
    const [tone, setTone] = useState("professional");
    const [content, setContent] = useState("");
    const [communicationId, setCommunicationId] = useState<string | null>(null);
    const [isGenerating, setIsGenerating] = useState(false);

    const candidates = Array.isArray(candidate) ? candidate : (candidate ? [candidate] : []);
    const isBulk = candidates.length > 1;

    const generateEmail = useGenerateEmail();

    const handleGenerate = async () => {
        if (candidates.length === 0) return;
        setIsGenerating(true);

        try {
            if (isBulk) {
                // In bulk mode, we generate one by one or we could add a bulk endpoint.
                // For now, let's just do it sequentially for simplicity or use the first one as template.
                const drafts: string[] = [];
                for (const c of candidates) {
                    const result = await generateEmail.mutateAsync({
                        candidateId: c.id,
                        type: type === "invitation" ? "interview_invite" : "rejection",
                        tone: tone as any,
                        locale,
                    });
                    drafts.push(`--- Candidate: ${c.name} ---\n${result.subject}\n\n${result.body}`);
                }
                setContent(drafts.join("\n\n" + "=".repeat(20) + "\n\n"));
                setCommunicationId(null); // Bulk doesn't have a single ID yet
            } else {
                const result = await generateEmail.mutateAsync({
                    candidateId: candidates[0].id,
                    type: type === "invitation" ? "interview_invite" : "rejection",
                    tone: tone as any,
                    locale,
                });
                setContent(`${result.subject}\n\n${result.body}`);
                setCommunicationId(result.id);
            }
        } catch (error) {
            // Error handled by hook toast
        } finally {
            setIsGenerating(false);
        }
    };

    const candidateIds = isBulk ? (candidate as CandidateCard[]).map(c => c.id).join(",") : (candidate as CandidateCard)?.id || "";

    useEffect(() => {
        if (isOpen && candidates.length > 0) {
            handleGenerate();
        }
    }, [isOpen, candidateIds, type, tone]);

    return {
        tone,
        setTone,
        content,
        setContent,
        communicationId,
        isGenerating,
        isBulk,
        candidates,
        handleGenerate,
    };
}
