"use client";

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal, FileText, Archive, Eye } from "lucide-react";

type JobStatus = "Active" | "Closed" | "Draft";

interface Job {
    id: string;
    title: string;
    department: string;
    createdDate: string;
    candidatesCount: number;
    status: JobStatus;
}

const MOCK_JOBS: Job[] = [
    {
        id: "1",
        title: "Senior Frontend Engineer",
        department: "Engineering",
        createdDate: "2023-10-15",
        candidatesCount: 45,
        status: "Active",
    },
    {
        id: "2",
        title: "Product Manager",
        department: "Product",
        createdDate: "2023-10-10",
        candidatesCount: 28,
        status: "Active",
    },
    {
        id: "3",
        title: "UX Designer",
        department: "Design",
        createdDate: "2023-09-28",
        candidatesCount: 12,
        status: "Draft",
    },
    {
        id: "4",
        title: "Backend Developer",
        department: "Engineering",
        createdDate: "2023-09-15",
        candidatesCount: 156,
        status: "Closed",
    },
    {
        id: "5",
        title: "Marketing Specialist",
        department: "Marketing",
        createdDate: "2023-10-01",
        candidatesCount: 8,
        status: "Active",
    },
];

const StatusBadge = ({ status }: { status: JobStatus }) => {
    switch (status) {
        case "Active":
            return (
                <Badge className="bg-green-100 text-green-700 hover:bg-green-100/80 border-green-200 shadow-none dark:bg-green-900/30 dark:text-green-400 dark:border-green-800">
                    Active
                </Badge>
            );
        case "Closed":
            return (
                <Badge variant="secondary" className="bg-slate-100 text-slate-700 hover:bg-slate-100/80 border-slate-200 shadow-none dark:bg-slate-800 dark:text-slate-400 dark:border-slate-800">
                    Closed
                </Badge>
            );
        case "Draft":
            return (
                <Badge variant="outline" className="text-slate-500 border-slate-300 dark:text-slate-400 dark:border-slate-700">
                    Draft
                </Badge>
            );
        default:
            return <Badge variant="outline">{status}</Badge>;
    }
};

export function JobsTable() {
    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Job Title</TableHead>
                        <TableHead>Department</TableHead>
                        <TableHead>Created</TableHead>
                        <TableHead>Candidates</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead className="text-right">Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {MOCK_JOBS.map((job) => (
                        <TableRow key={job.id}>
                            <TableCell className="font-medium">{job.title}</TableCell>
                            <TableCell>{job.department}</TableCell>
                            <TableCell>{job.createdDate}</TableCell>
                            <TableCell>{job.candidatesCount}</TableCell>
                            <TableCell>
                                <StatusBadge status={job.status} />
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
                                        <DropdownMenuItem>
                                            <Eye className="mr-2 h-4 w-4" />
                                            View Candidates
                                        </DropdownMenuItem>
                                        <DropdownMenuItem>
                                            <FileText className="mr-2 h-4 w-4" />
                                            Edit Job
                                        </DropdownMenuItem>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem className="text-red-600">
                                            <Archive className="mr-2 h-4 w-4" />
                                            Archive
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
}
