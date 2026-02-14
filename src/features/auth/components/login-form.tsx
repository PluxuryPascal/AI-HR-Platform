"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslations } from "next-intl";
import { Loader2 } from "lucide-react";

import { Link, useRouter } from "@/i18n/routing";
import { useAuth } from "@/store/use-auth";
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

const loginSchema = z.object({
    email: z.string().email(),
    password: z.string().min(6),
});

export function LoginForm() {
    const t = useTranslations("Auth");
    const router = useRouter();
    const { login } = useAuth();
    const [isLoading, setIsLoading] = useState(false);

    const form = useForm<z.infer<typeof loginSchema>>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    });

    async function onSubmit(values: z.infer<typeof loginSchema>) {
        setIsLoading(true);
        // Simulate network delay
        await new Promise((resolve) => setTimeout(resolve, 1000));

        login({
            id: "1",
            name: "John Doe",
            email: values.email,
            avatar: "https://github.com/shadcn.png",
        });

        setIsLoading(false);
        router.push("/dashboard");
    }

    // Common input class to ensure contrast on glass background
    const inputClass = "bg-white/50 dark:bg-black/50 border-white/20 dark:border-white/10";

    return (
        <div className="w-full bg-transparent">
            <div className="flex flex-col space-y-2 p-8 pb-4 text-center bg-transparent">
                <h2 className="text-2xl font-semibold tracking-tight">{t("signIn")}</h2>
                <p className="text-sm text-muted-foreground">{t("signInDesc")}</p>
            </div>

            <div className="p-8 pt-4 space-y-6 bg-transparent">
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 bg-transparent">
                        <FormField
                            control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem className="space-y-1">
                                    <FormLabel>{t("email")}</FormLabel>
                                    <FormControl>
                                        <Input placeholder="m@example.com" {...field} className={inputClass} />
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
                                        <PasswordInput {...field} className={inputClass} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <Button className="w-full bg-blue-600 hover:bg-blue-700 mt-2" type="submit" disabled={isLoading}>
                            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t("signInBtn")}
                        </Button>
                    </form>
                </Form>
            </div>

            <div className="flex justify-center p-8 pt-0 bg-transparent">
                <p className="text-sm text-muted-foreground p-0 font-normal">
                    {t("noAccount")}{" "}
                    <Link
                        href={{ pathname: '/auth', query: { mode: 'register' } }}
                        className="underline underline-offset-4 ml-1 font-medium text-primary hover:text-blue-600 transition-colors"
                        replace
                    >
                        {t("signUp")}
                    </Link>
                </p>
            </div>
        </div>
    );
}
