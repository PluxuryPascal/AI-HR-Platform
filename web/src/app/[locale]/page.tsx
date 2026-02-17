"use client";

import { GlobalDotGridBg } from "@/components/ui/global-dot-grid-bg";
import { LandingNavbar } from "@/components/layout/landing-navbar";
import { LandingHero } from "./_components/landing-hero";
import { LandingFeatures } from "./_components/landing-features";
import { motion, Variants } from "framer-motion";

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

export default function LandingPage() {
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

        <LandingNavbar />

        <main className="flex-1">
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            <LandingHero />
            <LandingFeatures />
          </motion.div>
        </main>
      </div>
    </>
  );
}
