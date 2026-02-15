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
    // Mock different profiles based on ID for testing
    if (id === "junior") {
        return {
            id,
            personal: {
                name: "Jordan Lee",
                email: "jordan.lee@example.com",
                phone: "+1 (555) 987-6543",
                location: "Austin, TX",
                role: "Junior Developer",
            },
            aiAnalysis: {
                score: 45,
                summary: "Enthusiastic learner with some bootcamp experience. Lacks commercial experience and depth in core concepts.",
                strengths: ["Fast learner", "Good communicator", "Team player"],
                weaknesses: ["No production experience", "Weak algorithm skills", "Unfamiliar with testing"],
                skills: ["HTML", "CSS", "Basic JavaScript"],
            },
            pdfUrl: "/dummy.pdf",
        };
    }

    if (id === "mid") {
        return {
            id,
            personal: {
                name: "Casey Taylor",
                email: "casey.taylor@example.com",
                phone: "+1 (555) 456-7890",
                location: "New York, NY",
                role: "Frontend Developer",
            },
            aiAnalysis: {
                score: 72,
                summary: "Solid contribution to previous projects. Good grasp of React but needs improvement in system design and performance optimization.",
                strengths: ["React Functional Components", "Redux Toolkit", "API Integration"],
                weaknesses: ["Performance profiling", "Advanced TypeScript types"],
                skills: ["React", "TypeScript", "Redux", "SASS"],
            },
            pdfUrl: "/dummy.pdf",
        };
    }

    // Default to Senior for 'senior' or any other ID
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
        pdfUrl: "/dummy.pdf",
    };
};
