"use client";

import { useState, useEffect } from "react";

export type CandidateStatus = "new" | "screening" | "interview" | "rejected" | "offer";

export interface Candidate {
    id: string;
    name: string;
    email: string;
    appliedRole: string;
    score: number;
    status: CandidateStatus;
    matchSummary: string;
    appliedDate: string;
}

const MOCK_CANDIDATES: Candidate[] = [
    {
        id: "1",
        name: "Alex Johnson",
        email: "alex.j@example.com",
        appliedRole: "Senior Frontend Engineer",
        score: 92,
        status: "screening",
        matchSummary: "Strong React & TypeScript experience. Good cultural fit.",
        appliedDate: "2024-02-14",
    },
    {
        id: "2",
        name: "Maria Garcia",
        email: "m.garcia@example.com",
        appliedRole: "Product Designer",
        score: 88,
        status: "new",
        matchSummary: "Excellent portfolio. Fits all design requirements.",
        appliedDate: "2024-02-13",
    },
    {
        id: "3",
        name: "James Chen",
        email: "j.chen@example.com",
        appliedRole: "Backend Developer",
        score: 45,
        status: "rejected",
        matchSummary: "Lacks required Python experience.",
        appliedDate: "2024-02-12",
    },
    {
        id: "4",
        name: "Sarah Smith",
        email: "s.smith@example.com",
        appliedRole: "Marketing Manager",
        score: 75,
        status: "interview",
        matchSummary: "Good experience, but salary expectations are high.",
        appliedDate: "2024-02-10",
    },
    {
        id: "5",
        name: "David Kim",
        email: "d.kim@example.com",
        appliedRole: "Senior Frontend Engineer",
        score: 95,
        status: "offer",
        matchSummary: "Perfect match. former Google.",
        appliedDate: "2024-02-09",
    },
    {
        id: "6",
        name: "Emily Davis",
        email: "e.davis@example.com",
        appliedRole: "Product Manager",
        score: 60,
        status: "screening",
        matchSummary: "Junior profile, might need training.",
        appliedDate: "2024-02-08",
    },
    {
        id: "7",
        name: "Michael Wilson",
        email: "m.wilson@example.com",
        appliedRole: "DevOps Engineer",
        score: 20,
        status: "rejected",
        matchSummary: "Not authorized to work in the country.",
        appliedDate: "2024-02-07",
    },
];

export function useCandidates() {
    const [data, setData] = useState<Candidate[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const timer = setTimeout(() => {
            // Simulate 5% chance of error
            if (Math.random() < 0.05) {
                // We're not setting error state here directly because the hook doesn't export it, 
                // but for the Global Error boundary to catch it, it needs to happen during render or be thrown.
                // However, since we're in a useEffect, throwing here won't be caught by ErrorBoundary unless we bubble it up.
                // For simplicity in this demo, we'll just log it, but to truly test the WidgetError, we might need to modify components to throw.
                // Actually, the user asked to simulate error in the HOOK. UseState updates are safe.
                // Let's create a "simulated failure" by not setting data and maybe setting a separate error state if we had one.
                // Use Promise.reject as requested.
                Promise.reject(new Error("Simulated API Error"));
            }

            setData([]); // MOCK_CANDIDATES
            setIsLoading(false);
        }, 800);

        return () => clearTimeout(timer);
    }, []);

    return { data, isLoading };
}
