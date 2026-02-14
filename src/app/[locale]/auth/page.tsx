"use client";

import { useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";

import { GlobalDotGridBg } from "@/components/ui/global-dot-grid-bg";
import { LoginForm } from "@/features/auth/components/login-form";
import { RegisterForm } from "@/features/auth/components/register-form";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { Link } from "@/i18n/routing";

export default function AuthPage() {
    const searchParams = useSearchParams();
    const mode = searchParams.get("mode");
    const [isLogin, setIsLogin] = useState(true);

    useEffect(() => {
        if (mode === "register") {
            setIsLogin(false);
        } else {
            setIsLogin(true);
        }
    }, [mode]);

    return (
        <div className="relative flex min-h-screen w-full flex-col items-center justify-center overflow-hidden bg-background selection:bg-blue-500/30">
            {/* Background with z-index -1 to sit behind content */}
            <div className="absolute inset-0 z-0">
                <GlobalDotGridBg />
            </div>

            {/* Back Button */}
            <div className="absolute top-4 left-4 z-20 md:top-8 md:left-8">
                <Button variant="ghost" size="sm" asChild className="gap-2 text-muted-foreground hover:text-foreground">
                    <Link href="/">
                        <ArrowLeft className="h-4 w-4" />
                        Back to Home
                    </Link>
                </Button>
            </div>

            <div className="relative z-10 w-full max-w-md px-4">
                {/* Animated Card Container - INLINE STYLES BYPASS TAILWIND */}
                <motion.div
                    layout
                    transition={{ type: "spring", stiffness: 300, damping: 30 }}
                    style={{
                        backdropFilter: "blur(4px)",
                        WebkitBackdropFilter: "blur(4px)"
                    }}
                    className="relative z-10 w-full max-w-md overflow-hidden rounded-[2.5rem] border border-white/50 bg-slate-900/5 shadow-[0_20px_50px_rgba(0,0,0,0.08)] dark:border-white/10 dark:bg-white/5 dark:shadow-2xl"
                >
                    {/* 3D Glare Overlay */}
                    <div className="pointer-events-none absolute inset-0 z-0 rounded-[2.5rem] bg-gradient-to-br from-white/50 via-white/10 to-transparent shadow-[inset_0_3px_20px_rgba(0,0,0,0.04),inset_0_1px_2px_rgba(255,255,255,1),inset_0_-1px_2px_rgba(0,0,0,0.05)] dark:from-white/15 dark:via-transparent dark:to-transparent dark:shadow-[inset_0_3px_20px_rgba(255,255,255,0.2),inset_0_1px_2px_rgba(255,255,255,0.8)]" />

                    <AnimatePresence mode="wait" initial={false}>
                        {isLogin ? (
                            <motion.div
                                key="login"
                                initial={{ opacity: 0, scale: 0.95, filter: "blur(10px)", y: 20 }}
                                animate={{ opacity: 1, scale: 1, filter: "blur(0px)", y: 0 }}
                                exit={{ opacity: 0, scale: 1.05, filter: "blur(10px)", y: -20 }}
                                transition={{ duration: 0.4, ease: "easeOut" }}
                                className="w-full bg-transparent relative z-10"
                            >
                                <LoginForm />
                            </motion.div>
                        ) : (
                            <motion.div
                                key="register"
                                initial={{ opacity: 0, scale: 0.95, filter: "blur(10px)", y: 20 }}
                                animate={{ opacity: 1, scale: 1, filter: "blur(0px)", y: 0 }}
                                exit={{ opacity: 0, scale: 1.05, filter: "blur(10px)", y: -20 }}
                                transition={{ duration: 0.4, ease: "easeOut" }}
                                className="w-full bg-transparent relative z-10"
                            >
                                <RegisterForm />
                            </motion.div>
                        )}
                    </AnimatePresence>
                </motion.div>

                {/* Footer Text */}
                <motion.p
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.5 }}
                    className="mt-6 text-center text-xs text-muted-foreground"
                >
                    By clicking continue, you agree to our{" "}
                    <Link href="#" className="underline underline-offset-4 hover:text-primary">
                        Terms of Service
                    </Link>{" "}
                    and{" "}
                    <Link href="#" className="underline underline-offset-4 hover:text-primary">
                        Privacy Policy
                    </Link>
                    .
                </motion.p>
            </div>
        </div>
    );
}
