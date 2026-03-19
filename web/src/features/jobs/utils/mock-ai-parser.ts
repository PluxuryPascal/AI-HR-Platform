import { JobFormValues, JobType } from "../schemas/job-schema";

export const parseJobDescription = async (text: string): Promise<Partial<JobFormValues>> => {
    // Simulate AI latency
    await new Promise((resolve) => setTimeout(resolve, 1500));

    // Simple keyword matching for "AI" simulation
    const lowerText = text.toLowerCase();

    let department: string = "Other";
    if (lowerText.includes("engineer") || lowerText.includes("developer") || lowerText.includes("code")) {
        department = "Engineering";
    } else if (lowerText.includes("design") || lowerText.includes("ui") || lowerText.includes("ux")) {
        department = "Design";
    } else if (lowerText.includes("product")) {
        department = "Product";
    } else if (lowerText.includes("market")) {
        department = "Marketing";
    } else if (lowerText.includes("sale")) {
        department = "Sales";
    }

    let type: typeof JobType[keyof typeof JobType] = JobType.Onsite;
    if (lowerText.includes("remote")) {
        type = JobType.Remote;
    } else if (lowerText.includes("hybrid")) {
        type = JobType.Hybrid;
    }

    // Extract a potential title (first line or first sentence)
    const lines = text.split('\n').filter(line => line.trim().length > 0);
    const potentialTitle = lines.length > 0 ? lines[0].substring(0, 100) : "Extracted Job Title";

    return {
        title: potentialTitle,
        department: department,
        description: text,
        requirements: ["Extracted requirement 1", "Extracted requirement 2", "Extracted requirement 3"], // Mock data
        type: type,
    };
};
