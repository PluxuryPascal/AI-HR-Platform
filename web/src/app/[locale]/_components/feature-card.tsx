"use client";

import { LucideIcon } from "lucide-react";
import { motion, Variants } from "framer-motion";
import {
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { SpotlightCard } from "@/components/ui/spotlight-card";
import { useTranslations } from "next-intl";

interface FeatureCardStyles {
    border: string;
    shadow: string;
    glow: string;
    iconBg: string;
    iconText: string;
    iconRing: string;
}

interface FeatureCardProps {
    icon: LucideIcon;
    titleKey: string;
    descKey: string;
    styles: FeatureCardStyles;
}

const itemVariants: Variants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
        opacity: 1,
        y: 0,
        transition: {
            duration: 0.5,
            ease: "easeOut",
        },
    },
};

export function FeatureCard({ icon: Icon, titleKey, descKey, styles }: FeatureCardProps) {
    const t = useTranslations("Landing");

    return (
        <motion.div variants={itemVariants} className="group">
            <SpotlightCard className={`h-full relative overflow-hidden rounded-2xl border border-slate-200/50 bg-white/50 p-8 backdrop-blur-md transition-all duration-300 hover:-translate-y-1 ${styles.border} ${styles.shadow} dark:border-slate-800/50 dark:bg-slate-900/50`}>
                <div className={`absolute -right-10 -top-10 h-32 w-32 rounded-full blur-3xl transition-all duration-500 ${styles.glow}`} />
                <CardHeader className="p-0 mb-6">
                    <div className={`flex h-12 w-12 items-center justify-center rounded-xl ${styles.iconBg} ${styles.iconText} ring-1 ${styles.iconRing}`}>
                        <Icon className="h-6 w-6" />
                    </div>
                </CardHeader>
                <CardContent className="p-0">
                    <CardTitle className="text-xl font-semibold text-slate-900 dark:text-slate-50 mb-3">
                        {t(titleKey)}
                    </CardTitle>
                    <CardDescription className="text-muted-foreground leading-relaxed text-sm">
                        {t(descKey)}
                    </CardDescription>
                </CardContent>
            </SpotlightCard>
        </motion.div>
    );
}
