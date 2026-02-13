export type Profile = {
    id: string;
    personal: {
        name: string;
        email: string;
        phone: string;
        location: string;
        role: string;
    };
    aiAnalysis: {
        score: number;
        summary: string;
        strengths: string[];
        weaknesses: string[];
        skills: string[];
    };
    pdfUrl: string;
};

export const getCandidateProfile = (id: string): Profile => {
    return {
        id,
        personal: {
            name: "Alex Morgan",
            email: "alex.morgan@example.com",
            phone: "+1 (555) 123-4567",
            location: "San Francisco, CA",
            role: "Senior Frontend Engineer",
        },
        aiAnalysis: {
            score: 88,
            summary: "Candidate shows strong proficiency in React ecosystem and modern frontend architecture. Previous experience leading small teams is a plus. Lacks deep backend knowledge but demonstrates quick learning ability.",
            strengths: [
                "Led a team of 5",
                "Optimized Web Vitals score by 40%",
                "Expert in React patterns",
                "Open source contributor",
            ],
            weaknesses: [
                "Lack of backend experience (Node.js/Python)",
                "Limited exposure to cloud infrastructure (AWS/Azure)",
            ],
            skills: ["React", "Next.js", "TypeScript", "Tailwind CSS", "Redux", "GraphQL"],
        },
        pdfUrl: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
    };
};
