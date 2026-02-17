import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

interface CandidateScoreBadgeProps {
    score: number;
    className?: string;
}

export function CandidateScoreBadge({ score, className }: CandidateScoreBadgeProps) {
    let colorClass = "bg-red-500/15 text-red-700 hover:bg-red-500/25 border-red-200";

    if (score >= 80) {
        colorClass = "bg-green-500/15 text-green-700 hover:bg-green-500/25 border-green-200";
    } else if (score >= 50) {
        colorClass = "bg-yellow-500/15 text-yellow-700 hover:bg-yellow-500/25 border-yellow-200";
    }

    return (
        <Badge variant="outline" className={cn("font-medium", colorClass, className)}>
            {score}% Match
        </Badge>
    );
}
