import { z } from "zod";

export const JobType = {
    Remote: "Remote",
    Onsite: "Onsite",
    Hybrid: "Hybrid",
} as const;

export const Department = {
    Engineering: "Engineering",
    Design: "Design",
    Product: "Product",
    Marketing: "Marketing",
    Sales: "Sales",
    HR: "HR",
    Finance: "Finance",
    Other: "Other",
} as const;

export const jobSchema = z.object({
    title: z.string().min(2, {
        message: "Title must be at least 2 characters.",
    }),
    department: z.nativeEnum(Department).or(z.string().min(1, "Department is required")),
    description: z.string().min(10, {
        message: "Description must be at least 10 characters.",
    }),
    requirements: z.array(z.string()).min(1, {
        message: "Add at least one requirement.",
    }),
    type: z.nativeEnum(JobType),
});

export type JobFormValues = z.infer<typeof jobSchema>;
