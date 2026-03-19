import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { 
    useGetAISettings, 
    useGetAIModels, 
    useUpdateAISetting,
    AISettings,
    AIModel 
} from "../api/use-ai-settings";

export type { AISettings, AIModel };

export function useAISettings() {
    const t = useTranslations("AISettings");
    
    const { data: settingsData, isLoading: isLoadingSettings } = useGetAISettings();
    const { data: modelsData, isLoading: isLoadingModels } = useGetAIModels();
    const { mutateAsync: updateSettingApi } = useUpdateAISetting();

    const isLoading = isLoadingSettings || isLoadingModels;

    const settings: AISettings = {
        team_id: settingsData?.team_id || "",
        api_key: settingsData?.api_key || "",
        parse_model: settingsData?.parse_model || "",
        score_model: settingsData?.score_model || "",
        embed_model: settingsData?.embed_model || "",
        chat_model: settingsData?.chat_model || "",
    };

    const models = modelsData || [];

    const updateSetting = async (key: keyof AISettings, value: string) => {
        try {
            await updateSettingApi({ field: key, value });
            toast.success(t("success"));
            return true;
        } catch (error) {
            toast.error(t("error"));
            return false;
        }
    };

    return {
        isLoading,
        models,
        settings,
        updateSetting,
    };
}
