import { cn } from "./utils";

/**
 * Reusable glassmorphism class strings.
 * Used across dashboard widgets and premium-styled cards.
 */

/** Dashboard card glass effect (overview chart, stats grid, recent activity) */
export const glassCard = cn(
    "border-white/10 bg-white/5 backdrop-blur-md",
    "shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)]",
    "hover:shadow-lg transition-all duration-300",
    "dark:bg-white/5 dark:border-white/10"
);

/** Floating toolbar glass pill (PDF viewer toolbar) */
export const glassToolbar = cn(
    "backdrop-blur-xl bg-white/70 dark:bg-zinc-900/80",
    "border border-white/20 dark:border-white/5",
    "shadow-2xl text-foreground",
    "transition-all duration-300"
);

/** Sticky header glass bar (dashboard header, landing navbar) */
export const stickyHeader = cn(
    "sticky top-0 z-50 w-full border-b",
    "bg-background/95 backdrop-blur",
    "supports-[backdrop-filter]:bg-background/60"
);
