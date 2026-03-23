"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslations } from "next-intl";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
    FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { useEffect } from "react";

const stageSchema = z.object({
    title: z.string().min(1, "Title is required"),
    is_terminal: z.boolean().default(false),
});

type StageFormValues = z.infer<typeof stageSchema>;

interface StageModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (values: StageFormValues) => void;
    initialValues?: Partial<StageFormValues>;
    title: string;
}

export function StageModal({
    isOpen,
    onClose,
    onSubmit,
    initialValues,
    title,
}: StageModalProps) {
    const t = useTranslations("Screening");

    const form = useForm<StageFormValues>({
        resolver: zodResolver(stageSchema) as any,
        defaultValues: {
            title: "",
            is_terminal: false,
            ...initialValues,
        },
    });

    useEffect(() => {
        if (isOpen) {
            form.reset({
                title: "",
                is_terminal: false,
                ...initialValues,
            });
        }
    }, [isOpen, initialValues, form]);

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                </DialogHeader>
                <Form {...(form as any)}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                            control={form.control as any}
                            name="title"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>{t("addStagePrompt")}</FormLabel>
                                    <FormControl>
                                        <Input placeholder="e.g. Technical Interview" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control as any}
                            name="is_terminal"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                        />
                                    </FormControl>
                                    <div className="space-y-1 leading-none">
                                        <FormLabel>
                                            {t("isTerminal")}
                                        </FormLabel>
                                    </div>
                                </FormItem>
                            )}
                        />
                        <DialogFooter>
                            <Button type="button" variant="outline" onClick={onClose}>
                                {t("bar.cancelSelection")}
                            </Button>
                            <Button type="submit">
                                {t("exitEditMode")}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
}
