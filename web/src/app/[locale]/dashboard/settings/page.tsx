
"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { LanguageSwitcher } from "@/components/shared/language-switcher";
import { useTranslations } from "next-intl";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { TeamManagement } from "@/features/settings/components/team-management";
import { AISettingsForm } from "@/features/settings/components/ai-settings-form";
import { useAuth } from "@/store/use-auth";
import { useGetProfile } from "@/features/profile/api/use-profile";

export default function SettingsPage() {
    const t = useTranslations("Settings");
    const { user } = useAuth();
    const { data: profile } = useGetProfile();
    const router = useRouter();
    const searchParams = useSearchParams();

    // Default to "general", but allow overriding via URL (though Tabs handles state internally usually)
    // We'll just rely on Tabs internal state or a URL param if we wanted deep linking.
    // For now, let's just stick to the requested redirect logic.

    const isOwner = profile?.role === "owner" || user?.role === "owner";

    // Handle tab change to enforce security on "team" tab access via UI
    const [activeTab, setActiveTab] = useState("general");

    useEffect(() => {
        // If query param 'tab' is 'team' and not owner, redirect to general (mock implementation)
        // Since Shadcn Tabs uses internal state by default, we just control the value.
    }, [user]);

    const handleTabChange = (value: string) => {
        if (value === "team" && !isOwner) {
            // Should not happen if button is hidden, but extra safety
            return;
        }
        setActiveTab(value);
    };

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between">
                <h2 className="text-3xl font-bold tracking-tight">{t("title")}</h2>
            </div>

            <Tabs value={activeTab} onValueChange={handleTabChange} className="space-y-4">
                <TabsList>
                    <TabsTrigger value="general">{t("general")}</TabsTrigger>
                    {isOwner && <TabsTrigger value="team">{t("team")}</TabsTrigger>}
                    <TabsTrigger value="ai">{t("ai")}</TabsTrigger>
                    <TabsTrigger value="billing">{t("billing")}</TabsTrigger>
                </TabsList>

                <TabsContent value="general" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>{t("general")}</CardTitle>
                            <CardDescription>
                                {t("languageDesc")}
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex flex-col space-y-2">
                                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                    {t("language")}
                                </label>
                                <LanguageSwitcher />
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {isOwner && (
                    <TabsContent value="team" className="space-y-4">
                        <TeamManagement />
                    </TabsContent>
                )}

                <TabsContent value="ai" className="space-y-4">
                    <AISettingsForm />
                </TabsContent>

                <TabsContent value="billing" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>{t("billing")}</CardTitle>
                            <CardDescription>
                                {t("billingDesc")}
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <p className="text-sm text-muted-foreground">{t("billingSoon")}</p>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
