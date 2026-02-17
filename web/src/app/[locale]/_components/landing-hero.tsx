"use client";

import { Link } from "@/i18n/routing";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ArrowRight } from "lucide-react";
import { motion, Variants } from "framer-motion";
import { useTranslations } from "next-intl";
import { useAuth } from "@/store/use-auth";

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

export function LandingHero() {
    const t = useTranslations("Landing");
    const { isAuthenticated } = useAuth();

    return (
        <section className="container mx-auto flex flex-col items-center justify-center px-4 py-24 text-center md:py-32 lg:py-40">
            <motion.div variants={itemVariants}>
                <Badge
                    variant="secondary"
                    className="mb-6 rounded-full px-4 py-1.5 text-sm font-medium text-blue-600 bg-blue-100/50 border-blue-200 dark:bg-blue-900/20 dark:text-blue-400 dark:border-blue-800"
                >
                    Powered by LLMs
                </Badge>
            </motion.div>

            <motion.h1
                variants={itemVariants}
                className="mb-6 max-w-4xl text-4xl font-extrabold tracking-tight text-slate-900 sm:text-5xl md:text-6xl lg:text-7xl dark:text-slate-50"
                dangerouslySetInnerHTML={{ __html: t.raw("heroTitle") }}
            />

            <motion.p
                variants={itemVariants}
                className="mb-10 max-w-2xl text-lg text-slate-600 md:text-xl dark:text-slate-400"
            >
                {t("heroDescription")}
            </motion.p>

            <motion.div
                variants={itemVariants}
                className="flex flex-col gap-4 sm:flex-row"
            >
                <Button size="lg" className="h-12 min-w-[160px] text-base gap-2 bg-blue-600 hover:bg-blue-500" asChild>
                    <Link href={isAuthenticated ? "/dashboard" : "/auth?mode=register"}>
                        {t("getStarted")} <ArrowRight className="h-4 w-4" />
                    </Link>
                </Button>
                <Button size="lg" variant="outline" className="h-12 min-w-[160px] text-base" asChild>
                    <Link href="/demo">{t("viewDemo")}</Link>
                </Button>
            </motion.div>
        </section>
    );
}
