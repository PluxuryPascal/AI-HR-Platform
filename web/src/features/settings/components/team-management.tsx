"use client";

import { useTranslations } from "next-intl";

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Plus } from "lucide-react";
import { toast } from "sonner";

import { InviteForm } from "./invite-form";
import { TeamTable } from "./team-table";
import { useInviteMember } from "../hooks/use-invite-member";

export function TeamManagement() {
    const t = useTranslations("Team");
    const { members, isInviteOpen, setIsInviteOpen, form, onSubmit } = useInviteMember();

    const handleSubmit = (values: Parameters<typeof onSubmit>[0]) => {
        onSubmit(values);
        toast.success(t("inviteModal.sendBtn") + " success!");
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
                            <InviteForm form={form} onSubmit={handleSubmit} />
                        </DialogContent>
                    </Dialog>
                </div>
            </CardHeader>
            <CardContent>
                <TeamTable members={members} />
            </CardContent>
        </Card>
    );
}
