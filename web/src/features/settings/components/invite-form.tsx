"use client";

import { useTranslations } from "next-intl";
import { Check, ChevronsUpDown } from "lucide-react";
import { UseFormReturn } from "react-hook-form";
import * as z from "zod";

import { cn } from "@/lib/utils";
import {
    DialogFooter,
} from "@/components/ui/dialog";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const mockJobs = [
    { id: "1", title: "Senior Frontend Developer" },
    { id: "2", title: "Product Manager" },
    { id: "3", title: "UI Designer" },
    { id: "4", title: "Backend Engineer" },
    { id: "5", title: "QA Specialist" },
];

export const inviteSchema = z.object({
    email: z.string().email({ message: "Invalid email address" }),
    role: z.enum(["owner", "recruiter", "manager"]),
    jobs: z.array(z.string()).optional(),
});

export type InviteFormValues = z.infer<typeof inviteSchema>;

interface InviteFormProps {
    form: UseFormReturn<InviteFormValues>;
    onSubmit: (values: InviteFormValues) => void;
}

export function InviteForm({ form, onSubmit }: InviteFormProps) {
    const t = useTranslations("Team");
    const watchRole = form.watch("role");

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t("inviteModal.emailLabel")}</FormLabel>
                            <FormControl>
                                <Input placeholder="colleague@company.com" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="role"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t("inviteModal.roleLabel")}</FormLabel>
                            <Select onValueChange={field.onChange} defaultValue={field.value}>
                                <FormControl>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select a role" />
                                    </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                    <SelectItem value="owner">{t("roles.owner")}</SelectItem>
                                    <SelectItem value="recruiter">{t("roles.recruiter")}</SelectItem>
                                    <SelectItem value="manager">{t("roles.manager")}</SelectItem>
                                </SelectContent>
                            </Select>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {watchRole === "manager" && (
                    <FormField
                        control={form.control}
                        name="jobs"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>{t("inviteModal.jobsLabel")}</FormLabel>
                                <Popover>
                                    <PopoverTrigger asChild>
                                        <FormControl>
                                            <Button
                                                variant="outline"
                                                role="combobox"
                                                className={cn(
                                                    "w-full justify-between",
                                                    !field.value?.length && "text-muted-foreground"
                                                )}
                                            >
                                                {field.value?.length
                                                    ? `${field.value.length} selected`
                                                    : t("inviteModal.searchJobs")}
                                                <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                            </Button>
                                        </FormControl>
                                    </PopoverTrigger>
                                    <PopoverContent className="w-[300px] p-0">
                                        <Command>
                                            <CommandInput placeholder={t("inviteModal.searchJobs")} />
                                            <CommandList>
                                                <CommandEmpty>No job found.</CommandEmpty>
                                                <CommandGroup>
                                                    {mockJobs.map((job) => (
                                                        <CommandItem
                                                            value={job.title}
                                                            key={job.id}
                                                            onSelect={() => {
                                                                const current = field.value || [];
                                                                const updated = current.includes(job.id)
                                                                    ? current.filter((id) => id !== job.id)
                                                                    : [...current, job.id];
                                                                field.onChange(updated);
                                                            }}
                                                        >
                                                            <Check
                                                                className={cn(
                                                                    "mr-2 h-4 w-4",
                                                                    field.value?.includes(job.id)
                                                                        ? "opacity-100"
                                                                        : "opacity-0"
                                                                )}
                                                            />
                                                            {job.title}
                                                        </CommandItem>
                                                    ))}
                                                </CommandGroup>
                                            </CommandList>
                                        </Command>
                                    </PopoverContent>
                                </Popover>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                )}

                <DialogFooter>
                    <Button type="submit">{t("inviteModal.sendBtn")}</Button>
                </DialogFooter>
            </form>
        </Form>
    );
}
