import { CandidateCard } from "../types";

type OutreachType = "rejection" | "invitation";
type OutreachTone = "professional" | "friendly" | "brief";

// Mock data for strengths and weaknesses since they are not yet in CandidateCard
const MOCK_ANALYSIS = {
    strengths: [
        "Strong React experience",
        "Good communication skills",
        "Experience with modern CI/CD",
        "Solid understanding of algorithms"
    ],
    weaknesses: [
        "Lack of deep Docker knowledge",
        "Limited experience with large-scale systems",
        "No experience with GraphQL",
        "Unfamiliar with Python"
    ]
};

const getRandomItem = (arr: string[]) => arr[Math.floor(Math.random() * arr.length)];

export function generateOutreachEmail(
    candidate: CandidateCard,
    type: OutreachType,
    tone: string,
    locale: string
): string {
    const isRejection = type === "rejection";
    // For now, we only have English templates fully implemented with dynamic content
    // In a real app, we'd have localized templates or use an LLM with locale context

    const weakness = getRandomItem(MOCK_ANALYSIS.weaknesses);
    const strength = getRandomItem(MOCK_ANALYSIS.strengths);

    if (tone === "brief") {
        if (isRejection) {
            return `Hi ${candidate.name},\n\nThank you for applying. After reviewing your profile, we have decided not to proceed. We are looking for someone with more ${weakness.toLowerCase().replace("lack of ", "")} experience.\n\nBest,\nHiring Team`;
        }
        return `Hi ${candidate.name},\n\nWe'd like to invite you for an interview! Your profile looks great. Let us know when you're available.\n\nBest,\nHiring Team`;
    }

    if (tone === "friendly") {
        if (isRejection) {
            return `Hi ${candidate.name} 👋,\n\nThanks so much for your interest in the ${candidate.role} position! We really appreciated getting to know your background.\n\nWhile we were impressed by your ${strength.toLowerCase()}, we're currently looking for a candidate with more depth in ${weakness.toLowerCase().replace("lack of ", "")}.\n\nWe'll keep your resume on file for future openings!\n\nWarmly,\nThe Hiring Team`;
        }
        return `Hi ${candidate.name}! 🚀,\n\nWe're super excited about your application for the ${candidate.role} role! Your experience with ${strength.toLowerCase()} really stood out to us.\n\nWe'd love to chat more in an interview. When would be a good time for you this week?\n\nCheers,\nThe Hiring Team`;
    }

    // Professional (Default)
    if (isRejection) {
        return `Dear ${candidate.name},\n\nThank you for giving us the opportunity to review your application for the ${candidate.role} position.\n\nAlthough your background is impressive, specifically your ${strength.toLowerCase()}, we have decided to move forward with other candidates who have stronger expertise in ${weakness.toLowerCase().replace("lack of ", "")}.\n\nWe wish you the best in your job search.\n\nSincerely,\nHiring Manager`;
    }

    return `Dear ${candidate.name},\n\nWe have reviewed your application for the ${candidate.role} position and would like to invite you to an interview.\n\nYour experience with ${strength.toLowerCase()} aligns well with what we are looking for. Please let us know your availability for a call in the coming days.\n\nBest regards,\nHiring Manager`;
}
