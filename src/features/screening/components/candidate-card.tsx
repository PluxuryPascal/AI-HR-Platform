"use client";

import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { CandidateCard as CandidateCardType } from "../types";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

interface CandidateCardProps {
    candidate: CandidateCardType;
    isOverlay?: boolean;
}

export function CandidateCard({ candidate, isOverlay }: CandidateCardProps) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({ id: candidate.id, data: { ...candidate } });

    const style = {
        transform: CSS.Translate.toString(transform),
        transition,
        opacity: isDragging ? 0.3 : 1,
    };

    if (isOverlay) {
        return (
            <Card className="cursor-grabbing shadow-xl rotate-2 scale-105 border-primary/50 relative z-50">
                <CardContent className="p-3 flex items-center space-x-4">
                    <Avatar>
                        <AvatarImage src={candidate.avatarUrl} alt={candidate.name} />
                        <AvatarFallback>{candidate.name.charAt(0)}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium leading-none">{candidate.name}</p>
                        <p className="text-xs text-muted-foreground">{candidate.role}</p>
                    </div>
                    <Badge
                        variant={candidate.score >= 80 ? "default" : "secondary"}
                        className={candidate.score >= 80 ? "bg-green-500" : ""}
                    >
                        {candidate.score}
                    </Badge>
                </CardContent>
            </Card>
        );
    }

    return (
        <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
            <Card
                className={`cursor-grab active:cursor-grabbing hover:shadow-md transition-shadow`}
            >
                <CardContent className="p-3 flex items-center space-x-4">
                    <Avatar>
                        <AvatarImage src={candidate.avatarUrl} alt={candidate.name} />
                        <AvatarFallback>{candidate.name.charAt(0)}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium leading-none">{candidate.name}</p>
                        <p className="text-xs text-muted-foreground">{candidate.role}</p>
                    </div>
                    <Badge
                        variant={candidate.score >= 80 ? "default" : "secondary"}
                        className={candidate.score >= 80 ? "bg-green-500 hover:bg-green-600" : ""}
                    >
                        {candidate.score}
                    </Badge>
                </CardContent>
            </Card>
        </div>
    );
}
