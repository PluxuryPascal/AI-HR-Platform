
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { LanguageSwitcher } from "@/components/shared/language-switcher";
import { useTranslations } from "next-intl";

export default function SettingsPage() {
    const t = useTranslations("Settings");

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <h2 className="text-3xl font-bold tracking-tight">{t("title")}</h2>

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
        </div>
    );
}
