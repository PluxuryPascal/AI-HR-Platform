import { useEffect, useState } from "react";
import { useLocale } from "next-intl";
import { CandidateCard } from "../types";
import { generateOutreachEmail } from "../utils/outreach-generator";

interface UseOutreachGeneratorOptions {
    isOpen: boolean;
    candidate: CandidateCard | CandidateCard[] | null;
    type: "rejection" | "invitation";
}

export function useOutreachGenerator({ isOpen, candidate, type }: UseOutreachGeneratorOptions) {
    const locale = useLocale();
    const [tone, setTone] = useState("professional");
    const [content, setContent] = useState("");
    const [isGenerating, setIsGenerating] = useState(false);

    const candidates = Array.isArray(candidate) ? candidate : (candidate ? [candidate] : []);
    const isBulk = candidates.length > 1;

    const handleGenerate = () => {
        if (candidates.length === 0) return;
        setIsGenerating(true);

        // Simulate AI delay
        setTimeout(() => {
            if (isBulk) {
                const drafts = candidates.map(c => {
                    return `--- Role: ${c.role} ---\n${generateOutreachEmail(c, type, tone, locale)}`;
                }).join("\n\n");

                // Note: We can't call t() in a hook, so we use a simpler prefix
                setContent(`Bulk outreach drafts:\n\n${drafts}`);
            } else {
                const newContent = generateOutreachEmail(candidates[0], type, tone, locale);
                setContent(newContent);
            }
            setIsGenerating(false);
        }, 600);
    };

    useEffect(() => {
        if (isOpen && candidates.length > 0) {
            handleGenerate();
        }
    }, [isOpen, candidate, type, tone]);

    return {
        tone,
        setTone,
        content,
        setContent,
        isGenerating,
        isBulk,
        candidates,
        handleGenerate,
    };
}
