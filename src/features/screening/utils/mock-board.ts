import { CandidateCard, ColumnId } from '../types';

export const initialColumns: Record<ColumnId, CandidateCard[]> = {
    new: [
        { id: 'c1', name: 'Alice Johnson', role: 'Frontend Dev', score: 85, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alice' },
        { id: 'c2', name: 'Bob Smith', role: 'Backend Dev', score: 72, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Bob' },
        { id: 'c3', name: 'Charlie Brown', role: 'Fullstack Dev', score: 90, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Charlie' },
        { id: 'c10', name: 'Jack Daniels', role: 'DevOps Engineer', score: 65, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Jack' },
    ],
    screening: [
        { id: 'c4', name: 'Diana Prince', role: 'Product Manager', score: 88, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Diana' },
        { id: 'c5', name: 'Evan Wright', role: 'Designer', score: 78, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Evan' },
        { id: 'c11', name: 'Karen Page', role: 'Legal Counsel', score: 91, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Karen' },
        { id: 'c12', name: 'Leo Leo', role: 'Marketing Lead', score: 75, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Leo' },
    ],
    interview: [
        { id: 'c6', name: 'Fiona Gallagher', role: 'HR Specialist', score: 92, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Fiona' },
        { id: 'c13', name: 'Mike Ross', role: 'Legal Associate', score: 89, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Mike' },
    ],
    offer: [
        { id: 'c7', name: 'George Martin', role: 'Tech Lead', score: 95, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=George' },
    ],
    rejected: [
        { id: 'c8', name: 'Hannah Abbott', role: 'Intern', score: 45, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Hannah' },
        { id: 'c9', name: 'Ian Malcolm', role: 'Chaos Theorist', score: 12, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Ian' },
        { id: 'c14', name: 'Noob Saibot', role: 'Shadow Ninja', score: 1, avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Noob' },
    ],
};
