"use client";

import { useTranslations } from "next-intl";
import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { PersonalInfoForm } from "@/features/profile/components/personal-info-form";
import { SecurityForm } from "@/features/profile/components/security-form";

export default function ProfilePage() {
    const t = useTranslations("Profile");

    return (
        <DashboardLayout>
            <div className="container mx-auto px-4 py-8 md:py-12 max-w-4xl">
                <div className="space-y-6">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
                        <p className="text-muted-foreground">{t("description")}</p>
                    </div>

                    <div className="grid gap-6">
                        <PersonalInfoForm />
                        <SecurityForm />
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
}
