"use client";

import { Link, usePathname } from "@/i18n/routing";
import { cn } from "@/lib/utils";
import { buttonVariants } from "@/components/ui/button";
import { dashboardConfig } from "@/config/nav";
import { useTranslations } from "next-intl";

export function SidebarNav({ className }: { className?: string }) {
    const pathname = usePathname();
    const t = useTranslations("Dashboard.nav");

    return (
        <nav className={cn("grid items-start gap-2", className)}>
            {dashboardConfig.map((item, index) => {
                const Icon = item.icon;
                const isActive = pathname === item.href;

                return (
                    <Link
                        key={index}
                        href={item.href}
                    >
                        <span
                            className={cn(
                                "group flex items-center rounded-md px-3 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground",
                                isActive ? "bg-accent text-accent-foreground" : "transparent"
                            )}
                        >
                            <Icon className="mr-2 h-4 w-4" />
                            <span>{t(item.key)}</span>
                        </span>
                    </Link>
                );
            })}
        </nav>
    );
}
