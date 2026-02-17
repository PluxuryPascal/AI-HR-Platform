import { CandidateCard } from "../types";

interface CsvExportData {
    candidates: CandidateCard[];
    aiData: Record<string, any>;
}

export function useCsvExport() {
    const handleExport = ({ candidates, aiData }: CsvExportData) => {
        const headers = ["Criteria", ...candidates.map(c => c.name)];
        const rows = [
            ["Role", ...candidates.map(c => c.role)],
            ["Score", ...candidates.map(c => c.score.toString())],
            ["Experience", ...candidates.map(c => aiData[c.id]?.experience || "")],
            ["Skills", ...candidates.map(c => aiData[c.id]?.skills.join(", ") || "")],
            ["Salary", ...candidates.map(c => aiData[c.id]?.salary || "")],
            ["Risks", ...candidates.map(c => aiData[c.id]?.risks || "")],
            ["Summary", ...candidates.map(c => aiData[c.id]?.summary || "")],
        ];

        const csvContent = [
            headers.join(","),
            ...rows.map(r => r.map(cell => `"${cell}"`).join(","))
        ].join("\n");

        const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);
        link.setAttribute("href", url);
        link.setAttribute("download", "candidate_comparison.csv");
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return { handleExport };
}
