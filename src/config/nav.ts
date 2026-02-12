import {
    LayoutDashboard,
    Briefcase,
    Users,
    FileText,
    Settings,
} from "lucide-react";

export const dashboardConfig = [
    {
        title: "Dashboard",
        href: "/dashboard",
        icon: LayoutDashboard,
    },
    {
        title: "Jobs",
        href: "/dashboard/jobs",
        icon: Briefcase,
    },
    {
        title: "Candidates",
        href: "/dashboard/candidates",
        icon: Users,
    },
    {
        title: "Screening",
        href: "/dashboard/screening",
        icon: FileText,
    },
    {
        title: "Settings",
        href: "/dashboard/settings",
        icon: Settings,
    },
];
