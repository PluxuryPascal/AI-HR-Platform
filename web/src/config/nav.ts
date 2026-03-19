import {
    LayoutDashboard,
    Briefcase,
    Users,
    FileText,
    Settings,
    Building2,
} from "lucide-react";

export const dashboardConfig = [
    {
        title: "Dashboard",
        href: "/dashboard",
        icon: LayoutDashboard,
        key: "dashboard",
    },
    {
        title: "Jobs",
        href: "/dashboard/jobs",
        icon: Briefcase,
        key: "jobs",
    },
    {
        title: "Candidates",
        href: "/dashboard/candidates",
        icon: Users,
        key: "candidates",
    },
    {
        title: "Screening",
        href: "/dashboard/screening",
        icon: FileText,
        key: "screening",
    },
    {
        title: "Departments",
        href: "/dashboard/departments",
        icon: Building2,
        key: "departments",
    },
    {
        title: "Settings",
        href: "/dashboard/settings",
        icon: Settings,
        key: "settings",
    },
];
