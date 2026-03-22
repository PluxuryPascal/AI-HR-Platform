import { z } from "zod";

export const JobType = {
    Remote: "Remote",
    Onsite: "Onsite",
    Hybrid: "Hybrid",
} as const;



export const jobSchema = z.object({
    title: z.string().min(2, {
        message: "Title must be at least 2 characters.",
    }),
    department: z.string().min(1, "Department is required"),
    description: z.string().min(10, {
        message: "Description must be at least 10 characters.",
    }),
    requirements: z.array(z.object({
        value: z.string().min(1, "Requirement cannot be empty")
    })).min(1, {
        message: "Add at least one requirement.",
    }),
    type: z.nativeEnum(JobType),
    salary_min: z.number().min(0, "Min salary must be 0 or more"),
    salary_max: z.number().min(0, "Max salary must be 0 or more"),
    currency: z.string().min(1, "Currency is required"),
});

export type JobFormValues = z.infer<typeof jobSchema>;
