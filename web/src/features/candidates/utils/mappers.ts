import { Candidate as BackendCandidate } from "../types/candidate";
import { CandidateCard } from "@/features/screening/types";

export function mapBackendCandidateToCard(c: BackendCandidate): CandidateCard {
    return {
        id: c.id,
        name: `${c.first_name || ""} ${c.last_name || ""}`.trim() || "Unknown",
        role: "Candidate",
        score: c.match_score || 0,
        avatarUrl: `https://api.dicebear.com/7.x/avataaars/svg?seed=${c.id}`,
        email: c.email || "-",
        appliedDate: new Date(c.created_at).toLocaleDateString(),
        status: c.stage_name || c.parsing_status || "New",
        matchSummary: "",
        scoreBreakdown: []
    };
}
