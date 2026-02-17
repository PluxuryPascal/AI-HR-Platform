"use client";

import { BrainCircuit, Wand2, MailCheck, KanbanSquare, Users, MoonStar } from "lucide-react";
import { FeatureCard } from "./feature-card";

const features = [
    {
        icon: BrainCircuit,
        titleKey: "explainableAI",
        descKey: "explainableAIDesc",
        styles: {
            border: "hover:border-blue-500/30",
            shadow: "hover:shadow-[0_0_40px_-15px_rgba(59,130,246,0.3)]",
            glow: "bg-blue-500/5 group-hover:bg-blue-500/10",
            iconBg: "bg-blue-500/10 dark:bg-blue-900/20",
            iconText: "text-blue-600 dark:text-blue-400",
            iconRing: "ring-blue-500/20 dark:ring-blue-400/20",
        },
    },
    {
        icon: Wand2,
        titleKey: "smartJobParsing",
        descKey: "smartJobParsingDesc",
        styles: {
            border: "hover:border-indigo-500/30",
            shadow: "hover:shadow-[0_0_40px_-15px_rgba(99,102,241,0.3)]",
            glow: "bg-indigo-500/5 group-hover:bg-indigo-500/10",
            iconBg: "bg-indigo-500/10 dark:bg-indigo-900/20",
            iconText: "text-indigo-600 dark:text-indigo-400",
            iconRing: "ring-indigo-500/20 dark:ring-indigo-400/20",
        },
    },
    {
        icon: MailCheck,
        titleKey: "contextualAIEmails",
        descKey: "contextualAIEmailsDesc",
        styles: {
            border: "hover:border-violet-500/30",
            shadow: "hover:shadow-[0_0_40px_-15px_rgba(139,92,246,0.3)]",
            glow: "bg-violet-500/5 group-hover:bg-violet-500/10",
            iconBg: "bg-violet-500/10 dark:bg-violet-900/20",
            iconText: "text-violet-600 dark:text-violet-400",
            iconRing: "ring-violet-500/20 dark:ring-violet-400/20",
        },
    },
    {
        icon: KanbanSquare,
        titleKey: "interactiveKanban",
        descKey: "interactiveKanbanDesc",
        styles: {
            border: "hover:border-blue-500/30",
            shadow: "hover:shadow-[0_0_40px_-15px_rgba(59,130,246,0.3)]",
            glow: "bg-blue-500/5 group-hover:bg-blue-500/10",
            iconBg: "bg-blue-500/10 dark:bg-blue-900/20",
            iconText: "text-blue-600 dark:text-blue-400",
            iconRing: "ring-blue-500/20 dark:ring-blue-400/20",
        },
    },
    {
        icon: Users,
        titleKey: "teamManagement",
        descKey: "teamManagementDesc",
        styles: {
            border: "hover:border-indigo-500/30",
            shadow: "hover:shadow-[0_0_40px_-15px_rgba(99,102,241,0.3)]",
            glow: "bg-indigo-500/5 group-hover:bg-indigo-500/10",
            iconBg: "bg-indigo-500/10 dark:bg-indigo-900/20",
            iconText: "text-indigo-600 dark:text-indigo-400",
            iconRing: "ring-indigo-500/20 dark:ring-indigo-400/20",
        },
    },
    {
        icon: MoonStar,
        titleKey: "smartPDFWorkspace",
        descKey: "smartPDFWorkspaceDesc",
        styles: {
            border: "hover:border-violet-500/30",
            shadow: "hover:shadow-[0_0_40px_-15px_rgba(139,92,246,0.3)]",
            glow: "bg-violet-500/5 group-hover:bg-violet-500/10",
            iconBg: "bg-violet-500/10 dark:bg-violet-900/20",
            iconText: "text-violet-600 dark:text-violet-400",
            iconRing: "ring-violet-500/20 dark:ring-violet-400/20",
        },
    },
];

export function LandingFeatures() {
    return (
        <section className="container mx-auto px-4 py-16 md:py-24">
            <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
                {features.map((feature, index) => (
                    <FeatureCard key={index} {...feature} />
                ))}
            </div>
        </section>
    );
}
