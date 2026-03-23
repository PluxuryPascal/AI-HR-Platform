import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

import { useAuth } from "@/store/use-auth";
import { useGetProfile, useUpdateProfile } from "../api/use-profile";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";

const personalInfoSchema = z.object({
    name: z.string().min(2),
    email: z.string().email(),
});

export function PersonalInfoForm() {
    const t = useTranslations("Profile.personalInfo");
    const { login } = useAuth();
    
    const { data: profile, isLoading: isFetching } = useGetProfile();
    const updateProfile = useUpdateProfile();

    const form = useForm<z.infer<typeof personalInfoSchema>>({
        resolver: zodResolver(personalInfoSchema),
        defaultValues: {
            name: "",
            email: "",
        },
    });

    useEffect(() => {
        if (profile) {
            form.reset({
                name: `${profile.first_name} ${profile.last_name}`.trim(),
                email: profile.email,
            });
        }
    }, [profile, form]);

    async function onSubmit(values: z.infer<typeof personalInfoSchema>) {
        const [firstName, ...lastNameParts] = values.name.trim().split(/\s+/);
        const lastName = lastNameParts.join(" ") || "-";

        updateProfile.mutate({
            first_name: firstName,
            last_name: lastName,
            email: values.email,
        }, {
            onSuccess: () => {
                toast.success(t("successMessage") || "Profile updated successfully");
                // Update local auth store
                if (profile) {
                    login({
                        ...profile,
                        firstName: firstName,
                        lastName: lastName,
                        email: values.email,
                    });
                }
            },
            onError: (error: any) => {
                toast.error(error?.response?.data?.message || "Failed to update profile");
            }
        });
    }

    const isLoading = isFetching || updateProfile.isPending;

    return (
        <Card>
            <CardHeader>
                <CardTitle>{t("title")}</CardTitle>
                <CardDescription>{t("desc")}</CardDescription>
            </CardHeader>
            <CardContent>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>{t("nameLabel")}</FormLabel>
                                    <FormControl>
                                        <Input placeholder="John Doe" disabled={isLoading} {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>{t("emailLabel")}</FormLabel>
                                    <FormControl>
                                        <Input placeholder="m@example.com" disabled={isLoading} {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <div className="flex justify-end">
                            <Button type="submit" disabled={isLoading}>
                                {updateProfile.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                {t("saveBtn")}
                            </Button>
                        </div>
                    </form>
                </Form>
            </CardContent>
        </Card>
    );
}
