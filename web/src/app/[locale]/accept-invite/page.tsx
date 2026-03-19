"use client";

import { useSearchParams } from "next/navigation";
import { motion } from "framer-motion";

import { GlobalDotGridBg } from "@/components/ui/global-dot-grid-bg";
import { AcceptInviteForm } from "@/features/auth/components/accept-invite-form";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { Link } from "@/i18n/routing";

export default function AcceptInvitePage() {
    const searchParams = useSearchParams();
    const token = searchParams.get("token");

    if (!token) {
        return (
            <div className="flex min-h-screen items-center justify-center">
                <div className="text-center space-y-4">
                    <h1 className="text-2xl font-bold">Invalid Invitation</h1>
                    <p className="text-muted-foreground">The invitation link is missing a token.</p>
                    <Button asChild>
                        <Link href="/">Back to Home</Link>
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className="relative flex min-h-screen w-full flex-col items-center justify-center overflow-hidden bg-slate-50 dark:bg-neutral-950 selection:bg-blue-500/30">
            {/* Background */}
            <div className="absolute inset-0 z-0 overflow-hidden pointer-events-none">
                <GlobalDotGridBg />
                <motion.div
                    initial={{ scale: 1, x: "-50%", y: "-50%" }}
                    animate={{ scale: [1, 1.1, 1] }}
                    transition={{ duration: 8, repeat: Infinity, repeatType: "mirror", ease: "easeInOut" }}
                    className="absolute top-[15%] left-[20%] h-[500px] w-[500px] rounded-full bg-indigo-500/40 blur-[120px] mix-blend-multiply dark:bg-indigo-500/20 dark:mix-blend-normal"
                />
                <motion.div
                    initial={{ scale: 1, x: "50%", y: "50%" }}
                    animate={{ scale: [1, 1.2, 1] }}
                    transition={{ duration: 10, repeat: Infinity, repeatType: "mirror", ease: "easeInOut" }}
                    className="absolute bottom-[15%] right-[20%] h-[400px] w-[400px] rounded-full bg-purple-500/40 blur-[120px] mix-blend-multiply dark:bg-purple-500/20 dark:mix-blend-normal"
                />
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
                <motion.div
                    layout
                    style={{
                        backgroundColor: 'rgba(255, 255, 255, 0.05)',
                        backdropFilter: 'blur(4px)',
                        WebkitBackdropFilter: 'blur(4px)',
                        boxShadow: 'inset 0 0 0 1px rgba(255, 255, 255, 0.1), 0 25px 50px -12px rgba(0, 0, 0, 0.25)'
                    }}
                    className="relative z-10 w-full max-w-md overflow-hidden rounded-[2.5rem] border border-white/20 dark:border-white/10 shadow-[0_8px_30px_rgb(0,0,0,0.04),0_20px_40px_rgba(0,0,0,0.04)]"
                >
                    <div className="pointer-events-none absolute inset-0 z-0 rounded-[2.5rem] bg-gradient-to-br from-white/50 via-white/10 to-transparent shadow-[inset_0_3px_20px_rgba(0,0,0,0.04),inset_0_1px_2px_rgba(255,255,255,1),inset_0_-1px_2px_rgba(0,0,0,0.05)] dark:from-white/15 dark:via-transparent dark:to-transparent dark:shadow-[inset_0_3px_20px_rgba(255,255,255,0.2),inset_0_1px_2px_rgba(255,255,255,0.8)]" />
                    
                    <div className="relative z-10">
                        <AcceptInviteForm token={token} />
                    </div>
                </motion.div>
            </div>
        </div>
    );
}
