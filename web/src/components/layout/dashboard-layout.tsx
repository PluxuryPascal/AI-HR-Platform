"use client";

import { SidebarNav } from "@/components/layout/sidebar-nav";
import { Header } from "@/components/layout/header";

export function DashboardLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="flex h-screen flex-col">
            {/* Header (sticky top) */}
            <Header />

            <div className="flex flex-1 overflow-hidden">
                {/* Sidebar (desktop only) */}
                <aside className="hidden h-full w-64 flex-col border-r bg-muted/40 md:flex">
                    <div className="flex-1 overflow-y-auto px-4 py-4">
                        <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
                            <SidebarNav />
                        </nav>
                    </div>
                </aside>

                {/* Main Content Area */}
                <main className="flex-1 overflow-y-auto bg-background p-4 lg:p-6">
                    {children}
                </main>
            </div>
        </div>
    );
}
