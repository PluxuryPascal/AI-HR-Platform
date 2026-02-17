"use client";

import { Link } from "@/i18n/routing";
import { Button } from "@/components/ui/button";
import { FileText, LayoutDashboard, LogOut } from "lucide-react";
import { ModeToggle } from "@/components/ui/mode-toggle";
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

export function LandingNavbar() {
    const t = useTranslations("Landing");
    const tAuth = useTranslations("Auth");
    const { isAuthenticated, user, logout } = useAuth();

    return (
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
    );
}
