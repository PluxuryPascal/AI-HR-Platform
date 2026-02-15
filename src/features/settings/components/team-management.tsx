"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/lib/utils";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
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
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Plus, MoreHorizontal } from "lucide-react";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { toast } from "sonner";

// Mock Data
type Role = "owner" | "recruiter" | "manager";
type Status = "active" | "pending";

interface Member {
    id: string;
    name: string;
    email: string;
    role: Role;
    status: Status;
    avatar?: string;
    assignedJobs?: string[]; // Array of Job IDs
}

const mockJobs = [
    { id: "1", title: "Senior Frontend Developer" },
    { id: "2", title: "Product Manager" },
    { id: "3", title: "UI Designer" },
    { id: "4", title: "Backend Engineer" },
    { id: "5", title: "QA Specialist" },
];

const initialMembers: Member[] = [
    {
        id: "1",
        name: "Alice Johnson",
        email: "alice@example.com",
        role: "owner",
        status: "active",
        avatar: "https://i.pravatar.cc/150?u=alice",
    },
    {
        id: "2",
        name: "Bob Smith",
        email: "bob@example.com",
        role: "recruiter",
        status: "active",
        avatar: "https://i.pravatar.cc/150?u=bob",
    },
    {
        id: "3",
        name: "Charlie Brown",
        email: "charlie@example.com",
        role: "manager",
        status: "pending",
        avatar: "https://i.pravatar.cc/150?u=charlie",
        assignedJobs: ["1", "3"],
    },
];

const inviteSchema = z.object({
    email: z.string().email({ message: "Invalid email address" }),
    role: z.enum(["owner", "recruiter", "manager"]),
    jobs: z.array(z.string()).optional(),
});

export function TeamManagement() {
    const t = useTranslations("Team");
    const [members, setMembers] = useState<Member[]>(initialMembers);
    const [isInviteOpen, setIsInviteOpen] = useState(false);

    const form = useForm<z.infer<typeof inviteSchema>>({
        resolver: zodResolver(inviteSchema),
        defaultValues: {
            email: "",
            role: "recruiter",
            jobs: [],
        },
    });

    const watchRole = form.watch("role");

    const onSubmit = (values: z.infer<typeof inviteSchema>) => {
        const newMember: Member = {
            id: Math.random().toString(36).substr(2, 9),
            name: values.email.split("@")[0], // Placeholder name
            email: values.email,
            role: values.role,
            status: "pending",
            assignedJobs: values.role === "manager" ? values.jobs : undefined,
        };

        setMembers([...members, newMember]);
        setIsInviteOpen(false);
        form.reset();
        toast.success(t("inviteModal.sendBtn") + " success!");
    };

    const getRoleBadgeColor = (role: Role) => {
        switch (role) {
            case "owner":
                return "default";
            case "recruiter":
                return "secondary";
            case "manager":
                return "outline";
            default:
                return "default";
        }
    };

    return (
        <Card>
            <CardHeader>
                <div className="flex items-center justify-between">
                    <div>
                        <CardTitle>{t("title")}</CardTitle>
                        <CardDescription>{t("description")}</CardDescription>
                    </div>
                    <Dialog open={isInviteOpen} onOpenChange={setIsInviteOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="mr-2 h-4 w-4" />
                                {t("inviteBtn")}
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>{t("inviteModal.title")}</DialogTitle>
                                <DialogDescription>
                                    {t("description")}
                                </DialogDescription>
                            </DialogHeader>
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
                        </DialogContent>
                    </Dialog>
                </div>
            </CardHeader>
            <CardContent>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{t("table.name")}</TableHead>
                            <TableHead>{t("table.role")}</TableHead>
                            <TableHead>Jobs</TableHead>
                            <TableHead>{t("table.status")}</TableHead>
                            <TableHead className="text-right">{t("table.actions")}</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {members.map((member) => (
                            <TableRow key={member.id}>
                                <TableCell className="font-medium">
                                    <div className="flex items-center gap-3">
                                        <Avatar className="h-9 w-9">
                                            <AvatarImage src={member.avatar} alt={member.name} />
                                            <AvatarFallback>{member.name.charAt(0).toUpperCase()}</AvatarFallback>
                                        </Avatar>
                                        <div className="flex flex-col">
                                            <span>{member.name}</span>
                                            <span className="text-xs text-muted-foreground">{member.email}</span>
                                        </div>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant={getRoleBadgeColor(member.role)}>
                                        {t(`roles.${member.role}`)}
                                    </Badge>
                                </TableCell>
                                <TableCell>
                                    {member.role === "manager" ? (
                                        <div className="flex flex-wrap gap-1">
                                            {member.assignedJobs?.length ? (
                                                member.assignedJobs.map(jobId => {
                                                    const job = mockJobs.find(j => j.id === jobId);
                                                    return job ? (
                                                        <Badge key={jobId} variant="outline" className="text-xs">
                                                            {job.title}
                                                        </Badge>
                                                    ) : null;
                                                })
                                            ) : (
                                                <span className="text-muted-foreground text-sm">-</span>
                                            )}
                                        </div>
                                    ) : (
                                        <span className="text-muted-foreground text-sm">All</span>
                                    )}
                                </TableCell>
                                <TableCell>
                                    <Badge variant={member.status === "active" ? "default" : "secondary"} className={member.status === "active" ? "bg-green-500 hover:bg-green-600" : "bg-yellow-500 hover:bg-yellow-600"}>
                                        {t(`status.${member.status}`)}
                                    </Badge>
                                </TableCell>
                                <TableCell className="text-right">
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button variant="ghost" className="h-8 w-8 p-0">
                                                <span className="sr-only">Open menu</span>
                                                <MoreHorizontal className="h-4 w-4" />
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                            <DropdownMenuItem onClick={() => navigator.clipboard.writeText(member.email)}>
                                                Copy Email
                                            </DropdownMenuItem>
                                            <DropdownMenuSeparator />
                                            <DropdownMenuItem className="text-red-600">Remove Member</DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
}
