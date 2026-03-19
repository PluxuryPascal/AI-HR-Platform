import { useState, useEffect } from "react";
import { toast } from "sonner";

export type AIModel = {
    id: string;
    name: string;
};

export type AISettings = {
    openrouter_api_key: string;
    parse_model: string;
    score_model: string;
    embed_model: string;
    chat_mode: string;
};

const MOCK_MODELS: AIModel[] = [
    { id: "openai/gpt-4o", name: "GPT-4o" },
    { id: "openai/gpt-4o-mini", name: "GPT-4o Mini" },
    { id: "anthropic/claude-3.5-sonnet", name: "Claude 3.5 Sonnet" },
    { id: "google/gemini-pro-1.5", name: "Gemini 1.5 Pro" },
    { id: "meta-llama/llama-3.1-70b-instruct", name: "Llama 3.1 70B" },
];

export function useAISettings() {
    const [isLoading, setIsLoading] = useState(false);
    const [models, setModels] = useState<AIModel[]>([]);
    const [settings, setSettings] = useState<AISettings>({
        openrouter_api_key: "",
        parse_model: "",
        score_model: "",
        embed_model: "",
        chat_mode: "",
    });

    useEffect(() => {
        // Simulate fetching models and settings
        const fetchData = async () => {
            setIsLoading(true);
            await new Promise((resolve) => setTimeout(resolve, 800));
            setModels(MOCK_MODELS);
            // Simulate initial settings being empty or some values
            setSettings({
                openrouter_api_key: "sk-or-v1-...",
                parse_model: "openai/gpt-4o-mini",
                score_model: "openai/gpt-4o",
                embed_model: "openai/gpt-4o-mini",
                chat_mode: "anthropic/claude-3.5-sonnet",
            });
            setIsLoading(false);
        };

        fetchData();
    }, []);

    const updateSetting = async (key: keyof AISettings, value: string) => {
        setIsLoading(true);
        try {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 1000));
            setSettings((prev) => ({ ...prev, [key]: value }));
            toast.success("Настройки успешно обновлены");
            return true;
        } catch (error) {
            toast.error("Ошибка при обновлении настроек");
            return false;
        } finally {
            setIsLoading(false);
        }
    };

    return {
        isLoading,
        models,
        settings,
        updateSetting,
    };
}
