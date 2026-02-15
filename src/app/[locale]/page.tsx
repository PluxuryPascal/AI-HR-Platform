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
import { FileText, Award, Zap, ArrowRight, User, LogOut, LayoutDashboard, BrainCircuit, Wand2, MailCheck, KanbanSquare, Users, MoonStar } from "lucide-react";
import { GlobalDotGridBg } from "@/components/ui/global-dot-grid-bg";
import { SpotlightCard } from "@/components/ui/spotlight-card";
import { ModeToggle } from "@/components/ui/mode-toggle";
import { motion, Variants } from "framer-motion";
import { useTranslations } from "next-intl";
import { useAuth } from "@/store/use-auth";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export default function LandingPage() {
  const t = useTranslations("Landing");
  const tAuth = useTranslations("Auth");
  const { isAuthenticated, user, logout } = useAuth();

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

      <div className="relative flex min-h-screen flex-col overflow-hidden selection:bg-blue-500/30">
        {/* Background Spheres */}
        <div className="absolute inset-0 z-0 overflow-hidden pointer-events-none">
          {/* Sphere 1 (Top Left) */}
          <motion.div
            initial={{ scale: 1, x: "-50%", y: "-50%" }}
            animate={{
              scale: [1, 1.1, 1],
            }}
            transition={{
              duration: 8,
              repeat: Infinity,
              repeatType: "mirror",
              ease: "easeInOut",
            }}
            className="absolute top-[15%] left-[20%] h-[500px] w-[500px] rounded-full bg-indigo-500/40 blur-[120px] mix-blend-multiply dark:bg-indigo-500/20 dark:mix-blend-normal"
          />

          {/* Sphere 2 (Bottom Right) */}
          <motion.div
            initial={{ scale: 1, x: "50%", y: "50%" }}
            animate={{
              scale: [1, 1.2, 1],
            }}
            transition={{
              duration: 10,
              repeat: Infinity,
              repeatType: "mirror",
              ease: "easeInOut",
            }}
            className="absolute bottom-[15%] right-[20%] h-[400px] w-[400px] rounded-full bg-purple-500/40 blur-[120px] mix-blend-multiply dark:bg-purple-500/20 dark:mix-blend-normal"
          />
        </div>
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

              {isAuthenticated ? (
                <div className="flex items-center gap-4">
                  <Button variant="ghost" asChild className="hidden sm:flex">
                    <Link href="/dashboard">
                      <LayoutDashboard className="mr-2 h-4 w-4" />
                      {t("dashboard")}
                    </Link>
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                        <Avatar className="h-8 w-8">
                          <AvatarImage src={user?.avatar} alt={user?.name} />
                          <AvatarFallback>{user?.name?.charAt(0)}</AvatarFallback>
                        </Avatar>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="w-56" align="end" forceMount>
                      <DropdownMenuLabel className="font-normal">
                        <div className="flex flex-col space-y-1">
                          <p className="text-sm font-medium leading-none">{user?.name}</p>
                          <p className="text-xs leading-none text-muted-foreground">
                            {user?.email}
                          </p>
                        </div>
                      </DropdownMenuLabel>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem asChild>
                        <Link href="/dashboard">
                          <LayoutDashboard className="mr-2 h-4 w-4" />
                          {t("dashboard")}
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={logout}>
                        <LogOut className="mr-2 h-4 w-4" />
                        <span>{tAuth("logout") || "Log out"}</span>
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              ) : (
                <>
                  <Button variant="ghost" asChild>
                    <Link href="/auth?mode=login">{tAuth("signIn")}</Link>
                  </Button>
                  <Button asChild className="bg-blue-600 hover:bg-blue-700">
                    <Link href="/auth?mode=register">{tAuth("signUp")}</Link>
                  </Button>
                </>
              )}
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
                  <Link href={isAuthenticated ? "/dashboard" : "/auth?mode=register"}>
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
                {[
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
                    }
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
                    }
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
                    }
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
                    }
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
                    }
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
                    }
                  },
                ].map((feature, index) => (
                  <motion.div key={index} variants={itemVariants} className="group">
                    <SpotlightCard className={`h-full relative overflow-hidden rounded-2xl border border-slate-200/50 bg-white/50 p-8 backdrop-blur-md transition-all duration-300 hover:-translate-y-1 ${feature.styles.border} ${feature.styles.shadow} dark:border-slate-800/50 dark:bg-slate-900/50`}>
                      <div className={`absolute -right-10 -top-10 h-32 w-32 rounded-full blur-3xl transition-all duration-500 ${feature.styles.glow}`} />
                      <CardHeader className="p-0 mb-6">
                        <div className={`flex h-12 w-12 items-center justify-center rounded-xl ${feature.styles.iconBg} ${feature.styles.iconText} ring-1 ${feature.styles.iconRing}`}>
                          <feature.icon className="h-6 w-6" />
                        </div>
                      </CardHeader>
                      <CardContent className="p-0">
                        <CardTitle className="text-xl font-semibold text-slate-900 dark:text-slate-50 mb-3">
                          {t(feature.titleKey)}
                        </CardTitle>
                        <CardDescription className="text-muted-foreground leading-relaxed text-sm">
                          {t(feature.descKey)}
                        </CardDescription>
                      </CardContent>
                    </SpotlightCard>
                  </motion.div>
                ))}
              </div>
            </section>
          </motion.div>
        </main>
      </div>
    </>
  );
}
