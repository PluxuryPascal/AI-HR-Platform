import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface Department {
    id: string;
    name: string;
}

interface DepartmentsState {
    departments: Department[];
    addDepartment: (name: string) => void;
    deleteDepartment: (id: string) => void;
}

const DEFAULT_DEPARTMENTS: Department[] = [
    { id: "dept-1", name: "Engineering" },
    { id: "dept-2", name: "Design" },
    { id: "dept-3", name: "Product" },
    { id: "dept-4", name: "Marketing" },
    { id: "dept-5", name: "Sales" },
    { id: "dept-6", name: "HR" },
    { id: "dept-7", name: "Finance" },
    { id: "dept-8", name: "Other" },
];

export const useDepartmentsStore = create<DepartmentsState>()(
    persist(
        (set) => ({
            departments: DEFAULT_DEPARTMENTS,
            addDepartment: (name: string) =>
                set((state) => ({
                    departments: [
                        ...state.departments,
                        { id: `dept-${Date.now()}`, name },
                    ],
                })),
            deleteDepartment: (id: string) =>
                set((state) => ({
                    departments: state.departments.filter((dept) => dept.id !== id),
                })),
        }),
        {
            name: "departments-storage",
        }
    )
);
