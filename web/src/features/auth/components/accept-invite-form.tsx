"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslations } from "next-intl";
import { Loader2 } from "lucide-react";

import { useRouter } from "@/i18n/routing";
import { useAcceptInvite } from "../api/use-auth";
import { AcceptInviteRequest } from "../types";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/ui/password-input";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";

const acceptInviteSchema = z.object({
    firstName: z.string().min(2).max(32),
    lastName: z.string().min(2).max(32),
    password: z.string().min(8).max(32),
});

interface AcceptInviteFormProps {
    token: string;
}

export function AcceptInviteForm({ token }: AcceptInviteFormProps) {
    const t = useTranslations("Auth");
    const router = useRouter();
    const acceptInvite = useAcceptInvite();

    const form = useForm<z.infer<typeof acceptInviteSchema>>({
        resolver: zodResolver(acceptInviteSchema),
        defaultValues: {
            firstName: "",
            lastName: "",
            password: "",
        },
    });

    async function onSubmit(values: z.infer<typeof acceptInviteSchema>) {
        const payload: AcceptInviteRequest = {
            token: token,
            password: values.password,
            first_name: values.firstName,
            last_name: values.lastName,
        };

        acceptInvite.mutate(payload, {
            onSuccess: () => {
                toast.success(t("successMessage") || "Account created successfully");
                router.push("/dashboard");
            },
            onError: (error: any) => {
                toast.error(error?.response?.data?.message || "Registration failed");
            }
        });
    }

    const inputClass = "bg-white/50 dark:bg-black/50 border-white/20 dark:border-white/10";
    const isLoading = acceptInvite.isPending;

    return (
        <div className="w-full bg-transparent">
            <div className="flex flex-col space-y-2 p-8 pb-4 text-center bg-transparent">
                <h2 className="text-2xl font-semibold tracking-tight">{t("signUp")}</h2>
                <p className="text-sm text-muted-foreground">{t("signUpDesc")}</p>
            </div>

            <div className="p-8 pt-4 space-y-6 bg-transparent">
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 bg-transparent">
                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="firstName"
                                render={({ field }) => (
                                    <FormItem className="space-y-1">
                                        <FormLabel>{t("firstName")}</FormLabel>
                                        <FormControl>
                                            <Input placeholder="John" disabled={isLoading} {...field} className={inputClass} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="lastName"
                                render={({ field }) => (
                                    <FormItem className="space-y-1">
                                        <FormLabel>{t("lastName")}</FormLabel>
                                        <FormControl>
                                            <Input placeholder="Doe" disabled={isLoading} {...field} className={inputClass} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        <FormField
                            control={form.control}
                            name="password"
                            render={({ field }) => (
                                <FormItem className="space-y-1">
                                    <FormLabel>{t("password")}</FormLabel>
                                    <FormControl>
                                        <PasswordInput disabled={isLoading} {...field} className={inputClass} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <Button className="w-full bg-blue-600 hover:bg-blue-700 mt-2" type="submit" disabled={isLoading}>
                            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t("signUpBtn")}
                        </Button>
                    </form>
                </Form>
            </div>
        </div>
    );
}
