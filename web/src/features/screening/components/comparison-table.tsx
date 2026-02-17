"use client";

import { useTranslations } from "next-intl";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Sparkles } from "lucide-react";
import { CandidateCard } from "../types";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

interface ComparisonTableProps {
    candidates: CandidateCard[];
    aiData: Record<string, any>;
}

export function ComparisonTable({ candidates, aiData }: ComparisonTableProps) {
    const tCriteria = useTranslations("Screening.criteria");

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead className="w-[200px]">Criteria</TableHead>
                    {candidates.map(candidate => (
                        <TableHead key={candidate.id} className="min-w-[250px]">
                            <div className="flex items-center gap-3">
                                <Avatar>
                                    <AvatarImage src={candidate.avatarUrl} />
                                    <AvatarFallback>{candidate.name[0]}</AvatarFallback>
                                </Avatar>
                                <div className="flex flex-col">
                                    <span className="font-medium text-foreground">{candidate.name}</span>
                                    <span className="text-xs font-normal">{candidate.role}</span>
                                </div>
                            </div>
                        </TableHead>
                    ))}
                </TableRow>
            </TableHeader>
            <TableBody>
                <TableRow>
                    <TableCell className="font-medium">Score</TableCell>
                    {candidates.map(c => (
                        <TableCell key={c.id}>
                            <Badge variant={c.score >= 80 ? "default" : "secondary"}>
                                {c.score}
                            </Badge>
                        </TableCell>
                    ))}
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">{tCriteria("experience")}</TableCell>
                    {candidates.map(c => (
                        <TableCell key={c.id}>{aiData[c.id]?.experience}</TableCell>
                    ))}
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">{tCriteria("skills")}</TableCell>
                    {candidates.map(c => (
                        <TableCell key={c.id}>
                            <div className="flex flex-wrap gap-1">
                                {aiData[c.id]?.skills.map((skill: string) => (
                                    <Badge key={skill} variant="outline" className="text-xs">
                                        {skill}
                                    </Badge>
                                ))}
                            </div>
                        </TableCell>
                    ))}
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">{tCriteria("salary")}</TableCell>
                    {candidates.map(c => (
                        <TableCell key={c.id}>{aiData[c.id]?.salary}</TableCell>
                    ))}
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium text-destructive">{tCriteria("risks")}</TableCell>
                    {candidates.map(c => (
                        <TableCell key={c.id} className="text-destructive/80">
                            {aiData[c.id]?.risks}
                        </TableCell>
                    ))}
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium text-primary">
                        <div className="flex items-center gap-2">
                            <Sparkles className="w-3 h-3" />
                            {tCriteria("summary")}
                        </div>
                    </TableCell>
                    {candidates.map(c => (
                        <TableCell key={c.id} className="italic text-muted-foreground">
                            &quot;{aiData[c.id]?.summary}&quot;
                        </TableCell>
                    ))}
                </TableRow>
            </TableBody>
        </Table>
    );
}
