import { create } from "zustand"
import { persist } from "zustand/middleware"
import { CandidateCard } from "@/features/screening/types"
import { initialColumns } from "@/features/screening/utils/mock-board"

export type JobStatus = "Active" | "Closed" | "Draft"

export interface Stage {
    id: string
    title: string
}

export interface Job {
    id: string
    title: string
    department: string
    createdDate: string
    candidatesCount: number
    status: JobStatus
    stages: Stage[]
}

interface JobsState {
    jobs: Job[]
    // Record<jobId, Record<stageId, candidates[]>>
    boards: Record<string, Record<string, CandidateCard[]>>
    addJob: (job: Omit<Job, "id" | "candidatesCount" | "createdDate" | "stages">) => void
    moveCandidate: (
        jobId: string,
        candidateId: string,
        sourceColumn: string,
        targetColumn: string,
        newIndex: number
    ) => void
    addStage: (jobId: string, title: string) => void
    removeStage: (jobId: string, stageId: string) => void
    renameStage: (jobId: string, stageId: string, newTitle: string) => void
    reorderStages: (jobId: string, stageIds: string[]) => void
}

const DEFAULT_STAGES: Stage[] = [
    { id: "new", title: "New Applications" },
    { id: "screening", title: "AI Screening" },
    { id: "interview", title: "Interview" },
    { id: "offer", title: "Offer Sent" },
    { id: "rejected", title: "Rejected" },
]

const MOCK_JOBS: Job[] = [
    {
        id: "1",
        title: "Senior Frontend Engineer",
        department: "Engineering",
        createdDate: "2024-03-01",
        candidatesCount: 12,
        status: "Active",
        stages: [...DEFAULT_STAGES],
    },
    {
        id: "2",
        title: "AI Researcher",
        department: "Engineering",
        createdDate: "2024-03-05",
        candidatesCount: 8,
        status: "Active",
        stages: [...DEFAULT_STAGES],
    },
    {
        id: "3",
        title: "Product Designer",
        department: "Design",
        createdDate: "2024-02-28",
        candidatesCount: 15,
        status: "Active",
        stages: [...DEFAULT_STAGES],
    },
]

// Initialize boards with candidates for each stage
const INITIAL_BOARDS: Record<string, Record<string, CandidateCard[]>> = {
    "1": initialColumns,
    "2": JSON.parse(JSON.stringify(initialColumns)),
    "3": JSON.parse(JSON.stringify(initialColumns)),
}

export const useJobsStore = create<JobsState>()(
    persist(
        (set) => ({
            jobs: MOCK_JOBS,
            boards: INITIAL_BOARDS,
            addJob: (job) =>
                set((state) => {
                    const newId = Math.random().toString(36).substring(7)
                    const newJob: Job = {
                        ...job,
                        id: newId,
                        candidatesCount: 0,
                        createdDate: new Date().toISOString().split("T")[0],
                        stages: [...DEFAULT_STAGES],
                    }
                    return {
                        jobs: [...state.jobs, newJob],
                        boards: {
                            ...state.boards,
                            [newId]: {
                                new: [],
                                screening: [],
                                interview: [],
                                offer: [],
                                rejected: [],
                            },
                        },
                    }
                }),
            moveCandidate: (jobId, candidateId, sourceColumn, targetColumn, newIndex) =>
                set((state) => {
                    const board = state.boards[jobId]
                    if (!board) return state

                    const sourceItems = [...(board[sourceColumn] || [])]
                    const targetItems =
                        sourceColumn === targetColumn
                            ? sourceItems
                            : [...(board[targetColumn] || [])]

                    const candidateIndex = sourceItems.findIndex((c) => c.id === candidateId)
                    if (candidateIndex === -1) return state

                    const [candidate] = sourceItems.splice(candidateIndex, 1)

                    if (sourceColumn === targetColumn) {
                        sourceItems.splice(newIndex, 0, candidate)
                        return {
                            boards: {
                                ...state.boards,
                                [jobId]: {
                                    ...board,
                                    [sourceColumn]: sourceItems,
                                },
                            },
                        }
                    } else {
                        targetItems.splice(newIndex, 0, candidate)
                        return {
                            boards: {
                                ...state.boards,
                                [jobId]: {
                                    ...board,
                                    [sourceColumn]: sourceItems,
                                    [targetColumn]: targetItems,
                                },
                            },
                        }
                    }
                }),
            addStage: (jobId, title) =>
                set((state) => {
                    const job = state.jobs.find((j) => j.id === jobId)
                    if (!job) return state

                    const stageId = Math.random().toString(36).substring(7)
                    const newStage: Stage = { id: stageId, title }

                    return {
                        jobs: state.jobs.map((j) =>
                            j.id === jobId ? { ...j, stages: [...j.stages, newStage] } : j
                        ),
                        boards: {
                            ...state.boards,
                            [jobId]: {
                                ...state.boards[jobId],
                                [stageId]: [],
                            },
                        },
                    }
                }),
            removeStage: (jobId, stageId) =>
                set((state) => {
                    const job = state.jobs.find((j) => j.id === jobId)
                    if (!job) return state

                    const candidatesToMove = state.boards[jobId]?.[stageId] || []
                    const otherStages = job.stages.filter((s) => s.id !== stageId)
                    const fallbackStageId = otherStages[0]?.id || "new"

                    const newBoard = { ...state.boards[jobId] }
                    delete newBoard[stageId]

                    // Move candidates to fallback stage (usually "new")
                    if (candidatesToMove.length > 0) {
                        newBoard[fallbackStageId] = [
                            ...(newBoard[fallbackStageId] || []),
                            ...candidatesToMove,
                        ]
                    }

                    return {
                        jobs: state.jobs.map((j) =>
                            j.id === jobId ? { ...j, stages: otherStages } : j
                        ),
                        boards: {
                            ...state.boards,
                            [jobId]: newBoard,
                        },
                    }
                }),
            renameStage: (jobId, stageId, newTitle) =>
                set((state) => ({
                    jobs: state.jobs.map((j) =>
                        j.id === jobId
                            ? {
                                  ...j,
                                  stages: j.stages.map((s) =>
                                      s.id === stageId ? { ...s, title: newTitle } : s
                                  ),
                              }
                            : j
                    ),
                })),
            reorderStages: (jobId, stageIds) =>
                set((state) => {
                    const job = state.jobs.find((j) => j.id === jobId)
                    if (!job) return state

                    const reorderedStages = stageIds
                        .map((id) => job.stages.find((s) => s.id === id))
                        .filter(Boolean) as Stage[]

                    return {
                        jobs: state.jobs.map((j) =>
                            j.id === jobId ? { ...j, stages: reorderedStages } : j
                        ),
                    }
                }),
        }),
        {
            name: "jobs-storage",
        }
    )
)
