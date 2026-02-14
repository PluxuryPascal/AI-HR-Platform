import { CandidateCard, ColumnId } from '../types';

export const initialColumns: Record<ColumnId, CandidateCard[]> = {
    new: [],
    screening: [],
    interview: [],
    offer: [],
    rejected: [],
};
// export const initialColumns: Record<ColumnId, CandidateCard[]> = {
//     new: [
//         { id: 'c1', name: 'Alice Johnson', role: 'Frontend Dev', score: 85, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alice' },
//         { id: 'c2', name: 'Bob Smith', role: 'Backend Dev', score: 72, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Bob' },
//     ],
//     screening: [
//         { id: 'c3', name: 'Charlie Brown', role: 'Fullstack Dev', score: 90, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Charlie' },
//         { id: 'c4', name: 'Diana Prince', role: 'Product Manager', score: 88, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Diana' },
//         { id: 'c5', name: 'Evan Wright', role: 'Designer', score: 78, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Evan' },
//     ],
//     interview: [
//         { id: 'c6', name: 'Fiona Gallagher', role: 'HR Specialist', score: 92, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Fiona' },
//     ],
//     offer: [
//         { id: 'c7', name: 'George Martin', role: 'Tech Lead', score: 95, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=George' },
//     ],
//     rejected: [
//         { id: 'c8', name: 'Hannah Abbott', role: 'Intern', score: 45, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Hannah' },
//         { id: 'c9', name: 'Ian Malcolm', role: 'chaos theorist', score: 12, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Ian' },
//     ],
// };
