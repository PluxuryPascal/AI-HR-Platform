"use client";

import { Link } from "@/i18n/routing";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { FileText, Award, Zap, ArrowRight } from "lucide-react";
import { GlobalDotGridBg } from "@/components/ui/global-dot-grid-bg";
import { SpotlightCard } from "@/components/ui/spotlight-card";
import { ModeToggle } from "@/components/ui/mode-toggle";
import { motion, Variants } from "framer-motion";
import { useTranslations } from "next-intl";

export default function LandingPage() {
  const t = useTranslations("Landing");

  const containerVariants: Variants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
        delayChildren: 0.2,
      },
    },
  };

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

  return (
    <>
      <GlobalDotGridBg />

      <div className="relative flex min-h-screen flex-col">
        {/* Navbar */}
        <header className="sticky top-0 z-50 w-full border-b border-slate-200 bg-white/50 backdrop-blur supports-[backdrop-filter]:bg-white/40 dark:border-slate-800 dark:bg-black/80">
          <div className="container mx-auto flex h-16 items-center justify-between px-4 md:px-6">
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white">
                <FileText className="h-5 w-5" />
              </div>
              <span className="text-xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
                ResumeAI
              </span>
            </div>
            <nav className="flex items-center gap-4">
              <ModeToggle />
              <Button variant="ghost" asChild>
                <Link href="/login">Login</Link>
              </Button>
              <Button asChild className="bg-blue-600 hover:bg-blue-700">
                <Link href="/dashboard">{t("getStarted")}</Link>
              </Button>
            </nav>
          </div>
        </header>

        <main className="flex-1">
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            {/* Hero Section */}
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
                  <Link href="/dashboard">
                    {t("getStarted")} <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
                <Button size="lg" variant="outline" className="h-12 min-w-[160px] text-base" asChild>
                  <Link href="/demo">{t("viewDemo")}</Link>
                </Button>
              </motion.div>
            </section>

            {/* Features Section */}
            <section className="container mx-auto px-4 py-16 md:py-24">
              <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
                {/* Feature 1 */}
                <motion.div variants={itemVariants}>
                  <SpotlightCard className="h-full">
                    <CardHeader>
                      <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400">
                        <FileText className="h-6 w-6" />
                      </div>
                      <CardTitle className="text-xl">{t("smartParsing")}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <CardDescription className="text-base">
                        {t("smartParsingDesc")}
                      </CardDescription>
                    </CardContent>
                  </SpotlightCard>
                </motion.div>

                {/* Feature 2 */}
                <motion.div variants={itemVariants}>
                  <SpotlightCard className="h-full">
                    <CardHeader>
                      <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-indigo-100 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400">
                        <Award className="h-6 w-6" />
                      </div>
                      <CardTitle className="text-xl">{t("biasFree")}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <CardDescription className="text-base">
                        {t("biasFreeDesc")}
                      </CardDescription>
                    </CardContent>
                  </SpotlightCard>
                </motion.div>

                {/* Feature 3 */}
                <motion.div variants={itemVariants}>
                  <SpotlightCard className="h-full">
                    <CardHeader>
                      <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-violet-100 text-violet-600 dark:bg-violet-900/30 dark:text-violet-400">
                        <Zap className="h-6 w-6" />
                      </div>
                      <CardTitle className="text-xl">{t("instantInsights")}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <CardDescription className="text-base">
                        {t("instantInsightsDesc")}
                      </CardDescription>
                    </CardContent>
                  </SpotlightCard>
                </motion.div>
              </div>
            </section>
          </motion.div>
        </main>

        {/* Footer */}
        <footer className="border-t border-slate-200 bg-white/50 backdrop-blur-sm py-12 dark:border-slate-800 dark:bg-black/80">
          <div className="container mx-auto flex flex-col items-center justify-between gap-4 px-4 md:flex-row text-center md:text-left">
            <p className="text-sm text-slate-500 dark:text-slate-400">
              © {new Date().getFullYear()} ResumeAI. All rights reserved.
            </p>
            <div className="flex gap-6 text-sm text-slate-500 dark:text-slate-400">
              <Link href="#" className="hover:text-blue-600 hover:underline">
                Privacy Policy
              </Link>
              <Link href="#" className="hover:text-blue-600 hover:underline">
                Terms of Service
              </Link>
              <Link href="#" className="hover:text-blue-600 hover:underline">
                Contact
              </Link>
            </div>
          </div>
        </footer>
      </div>
    </>
  );
}
