"use client";

import { useAISettings, AISettings, AIModel } from "../hooks/use-ai-settings";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslations } from "next-intl";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, Save } from "lucide-react";
import { useEffect, useState } from "react";

function APIKeyForm({ initialValue, onSave }: { initialValue: string, onSave: (val: string) => Promise<boolean> }) {
    const t = useTranslations("AISettings");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const apiKeySchema = z.object({
        api_key: z.string().min(1, t("error")),
    });

    const form = useForm({
        resolver: zodResolver(apiKeySchema),
        defaultValues: { api_key: initialValue },
    });

    useEffect(() => {
        form.reset({ api_key: initialValue });
    }, [initialValue, form]);

    const onSubmit = async (values: z.infer<typeof apiKeySchema>) => {
        setIsSubmitting(true);
        await onSave(values.api_key);
        setIsSubmitting(false);
        form.reset(values);
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>{t("apiKeyTitle")}</CardTitle>
                <CardDescription>{t("apiKeyDesc")}</CardDescription>
            </CardHeader>
            <CardContent>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="flex items-end gap-4">
                        <FormField
                            control={form.control}
                            name="api_key"
                            render={({ field }) => (
                                <FormItem className="flex-1">
                                    <FormLabel>API Key</FormLabel>
                                    <FormControl>
                                        <Input type="password" placeholder="sk-or-v1-..." {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <Button type="submit" disabled={isSubmitting || !form.formState.isDirty}>
                            {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <Save className="w-4 h-4 mr-2" />}
                            {t("save")}
                        </Button>
                    </form>
                </Form>
            </CardContent>
        </Card>
    );
}

function ModelSelectionForm({ 
    title, 
    description, 
    label, 
    fieldName, 
    initialValue, 
    models, 
    onSave 
}: { 
    title: string;
    description: string;
    label: string;
    fieldName: keyof AISettings;
    initialValue: string;
    models: AIModel[];
    onSave: (val: string) => Promise<boolean>;
}) {
    const t = useTranslations("AISettings");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const schema = z.object({ [fieldName]: z.string().min(1, "Model is required") });
    const form = useForm({
        resolver: zodResolver(schema),
        defaultValues: { [fieldName]: initialValue },
    });

    useEffect(() => {
        form.reset({ [fieldName]: initialValue });
    }, [initialValue, form, fieldName]);

    const onSubmit = async (values: any) => {
        setIsSubmitting(true);
        await onSave(values[fieldName]);
        setIsSubmitting(false);
        form.reset(values);
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>{title}</CardTitle>
                <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="flex items-end gap-4">
                        <FormField
                            control={form.control}
                            name={fieldName as any}
                            render={({ field }) => (
                                <FormItem className="flex-1">
                                    <FormLabel>{label}</FormLabel>
                                    <Select onValueChange={field.onChange} value={field.value}>
                                        <FormControl>
                                            <SelectTrigger>
                                                <SelectValue placeholder="Select a model" />
                                            </SelectTrigger>
                                        </FormControl>
                                        <SelectContent>
                                            {models.map((model) => (
                                                <SelectItem key={model.id} value={model.id}>
                                                    {model.name}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <Button type="submit" disabled={isSubmitting || !form.formState.isDirty}>
                            {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <Save className="w-4 h-4 mr-2" />}
                            {t("save")}
                        </Button>
                    </form>
                </Form>
            </CardContent>
        </Card>
    );
}

export function AISettingsForm() {
    const t = useTranslations("AISettings");
    const { settings, models, updateSetting, isLoading } = useAISettings();

    if (isLoading && models.length === 0) {
        return (
            <div className="flex items-center justify-center p-12">
                <Loader2 className="w-8 h-8 animate-spin text-primary" />
            </div>
        );
    }

    const textModels = models.filter(m => !m.is_embedding);
    const embedModels = models.filter(m => m.is_embedding);

    return (
        <div className="space-y-6">
            <APIKeyForm 
                initialValue={settings.api_key || ""} 
                onSave={(val) => updateSetting("api_key", val)} 
            />
            
            <ModelSelectionForm 
                title={t("parseModelTitle")}
                description={t("parseModelDesc")}
                label={t("parseModelTitle")}
                fieldName="parse_model"
                initialValue={settings.parse_model || ""}
                models={textModels}
                onSave={(val) => updateSetting("parse_model", val)}
            />

            <ModelSelectionForm 
                title={t("scoreModelTitle")}
                description={t("scoreModelDesc")}
                label={t("scoreModelTitle")}
                fieldName="score_model"
                initialValue={settings.score_model || ""}
                models={textModels}
                onSave={(val) => updateSetting("score_model", val)}
            />

            <ModelSelectionForm 
                title={t("embedModelTitle")}
                description={t("embedModelDesc")}
                label={t("embedModelTitle")}
                fieldName="embed_model"
                initialValue={settings.embed_model || ""}
                models={embedModels}
                onSave={(val) => updateSetting("embed_model", val)}
            />

            <ModelSelectionForm 
                title={t("chatModeTitle")}
                description={t("chatModeDesc")}
                label={t("chatModeTitle")}
                fieldName="chat_model"
                initialValue={settings.chat_model || ""}
                models={textModels}
                onSave={(val) => updateSetting("chat_model", val)}
            />
        </div>
    );
}
