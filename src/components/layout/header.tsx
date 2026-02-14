"use client";

import { Link, usePathname } from "@/i18n/routing";
import { MobileNav } from "@/components/layout/mobile-nav";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { LogOut, User } from "lucide-react";
import { ModeToggle } from "@/components/ui/mode-toggle";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import React from "react";
import { useTranslations } from "next-intl";
import { dashboardConfig } from "@/config/nav";
import { useAuth } from "@/store/use-auth";

export function Header() {
    const pathname = usePathname();
    const segments = pathname.split("/").filter((item) => item !== "");
    const t = useTranslations("Dashboard.header");
    const tNav = useTranslations("Dashboard.nav");
    const { user, logout } = useAuth();

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="flex h-14 items-center px-4 md:px-6">
                <div className="mr-4 hidden md:flex">
                    <Link className="mr-6 flex items-center space-x-2 font-bold" href="/">
                        <span className="hidden sm:inline-block">Resume Screener</span>
                    </Link>
                </div>
                <MobileNav />

                <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
                    <div className="w-full flex-1 md:w-auto md:flex-none">
                        <Breadcrumb>
                            <BreadcrumbList>
                                <BreadcrumbItem>
                                    <BreadcrumbLink asChild>
                                        <Link href="/dashboard">{tNav("dashboard")}</Link>
                                    </BreadcrumbLink>
                                </BreadcrumbItem>
                                {segments.length > 1 &&
                                    segments.slice(1).map((segment, index) => {
                                        const href = `/dashboard/${segments.slice(1, index + 2).join("/")}`;
                                        const isLast = index === segments.length - 2;

                                        // Check if the segment corresponds to a known navigation item
                                        const navItem = dashboardConfig.find(item => item.key === segment);
                                        const title = navItem
                                            ? tNav(navItem.key)
                                            : (segment === "new" ? tNav("new") : segment.charAt(0).toUpperCase() + segment.slice(1));

                                        return (
                                            <React.Fragment key={href}>
                                                <BreadcrumbSeparator />
                                                <BreadcrumbItem>
                                                    {isLast ? (
                                                        <BreadcrumbPage>{title}</BreadcrumbPage>
                                                    ) : (
                                                        <BreadcrumbLink asChild>
                                                            <Link href={href}>{title}</Link>
                                                        </BreadcrumbLink>
                                                    )}
                                                </BreadcrumbItem>
                                            </React.Fragment>
                                        );
                                    })}
                            </BreadcrumbList>
                        </Breadcrumb>
                    </div>

                    <div className="flex items-center gap-2">
                        <ModeToggle />
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                                    <Avatar className="h-8 w-8">
                                        <AvatarImage src={user?.avatar} alt={user?.name || "User"} />
                                        <AvatarFallback>{user?.name?.charAt(0) || "U"}</AvatarFallback>
                                    </Avatar>
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent className="w-56" align="end" forceMount>
                                <DropdownMenuLabel className="font-normal">
                                    <div className="flex flex-col space-y-1">
                                        <p className="text-sm font-medium leading-none">{user?.name || "User"}</p>
                                        <p className="text-xs leading-none text-muted-foreground">
                                            {user?.email || "user@example.com"}
                                        </p>
                                    </div>
                                </DropdownMenuLabel>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>
                                    <User className="mr-2 h-4 w-4" />
                                    <span>{t("profile")}</span>
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={logout}>
                                    <LogOut className="mr-2 h-4 w-4" />
                                    <span>{t("logout")}</span>
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                </div>
            </div>
        </header>
    );
}

