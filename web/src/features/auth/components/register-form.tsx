"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslations } from "next-intl";
import { Loader2 } from "lucide-react";

import { Link, useRouter } from "@/i18n/routing";
import { useRegister } from "../api/use-auth";
import { RegisterRequest } from "../types";
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

const registerSchema = z.object({
    firstName: z.string().min(2).max(32),
    lastName: z.string().min(2).max(32),
    email: z.string().email(),
    password: z.string().min(8).max(32),
    teamName: z.string().min(3).max(32),
});

export function RegisterForm() {
    const t = useTranslations("Auth");
    const router = useRouter();
    const register = useRegister();

    const form = useForm<z.infer<typeof registerSchema>>({
        resolver: zodResolver(registerSchema),
        defaultValues: {
            firstName: "",
            lastName: "",
            email: "",
            password: "",
            teamName: "",
        },
    });

    async function onSubmit(values: z.infer<typeof registerSchema>) {
        const payload: RegisterRequest = {
            email: values.email,
            password: values.password,
            first_name: values.firstName,
            last_name: values.lastName,
            team_name: values.teamName,
        };

        register.mutate(payload, {
            onSuccess: () => {
                toast.success(t("successMessage") || "Account created successfully");
                router.push("/dashboard");
            },
            onError: (error: any) => {
                toast.error(error?.response?.data?.message || "Registration failed");
            }
        });
    }

    // Common input class to ensure contrast on glass background
    const inputClass = "bg-white/50 dark:bg-black/50 border-white/20 dark:border-white/10";

    const isLoading = register.isPending;

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
                            name="email"
                            render={({ field }) => (
                                <FormItem className="space-y-1">
                                    <FormLabel>{t("email")}</FormLabel>
                                    <FormControl>
                                        <Input placeholder="m@example.com" disabled={isLoading} {...field} className={inputClass} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="teamName"
                            render={({ field }) => (
                                <FormItem className="space-y-1">
                                    <FormLabel>{t("teamName")}</FormLabel>
                                    <FormControl>
                                        <Input placeholder="Acme Inc." disabled={isLoading} {...field} className={inputClass} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
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

            <div className="flex justify-center p-8 pt-0 bg-transparent">
                <p className="text-sm text-muted-foreground p-0 font-normal">
                    {t("hasAccount")}{" "}
                    <Link
                        href={{ pathname: '/auth', query: { mode: 'login' } }}
                        className="underline underline-offset-4 ml-1 font-medium text-primary hover:text-blue-600 transition-colors"
                        replace
                    >
                        {t("signIn")}
                    </Link>
                </p>
            </div>
        </div>
    );
}
