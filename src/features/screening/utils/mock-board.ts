import { CandidateCard, ColumnId } from '../types';

export const initialColumns: Record<ColumnId, CandidateCard[]> = {
    new: [
        {
            id: 'c1',
            name: 'Alice Johnson',
            role: 'Frontend Dev',
            score: 85,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alice',
            email: 'alice.j@example.com',
            appliedDate: '2024-02-14',
            matchSummary: 'Strong React skills, good cultural fit.',
            scoreBreakdown: [
                { id: '1', type: 'positive', text: 'Senior React experience', impact: '+20%' },
                { id: '2', type: 'neutral', text: 'Remote only', impact: '0%' }
            ]
        },
        {
            id: 'c2',
            name: 'Bob Smith',
            role: 'Backend Dev',
            score: 72,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Bob',
            email: 'bob.smith@example.com',
            appliedDate: '2024-02-13',
            matchSummary: 'Good SQL knowledge, lacks some Python experience.',
            scoreBreakdown: [
                { id: '1', type: 'negative', text: 'Low Python experience', impact: '-10%' },
                { id: '2', type: 'positive', text: 'Database optimization', impact: '+15%' }
            ]
        },
        {
            id: 'c3',
            name: 'Charlie Brown',
            role: 'Fullstack Dev',
            score: 90,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Charlie',
            email: 'charlie.b@example.com',
            appliedDate: '2024-02-12',
            matchSummary: 'Versatile developer with strong portfolio.',
            scoreBreakdown: [
                { id: '1', type: 'positive', text: 'Fullstack versatility', impact: '+25%' }
            ]
        },
        {
            id: 'c10',
            name: 'Jack Daniels',
            role: 'DevOps Engineer',
            score: 65,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Jack',
            email: 'jack.d@example.com',
            appliedDate: '2024-02-10',
            matchSummary: 'Intermediate DevOps skills, learning fast.',
        },
    ],
    screening: [
        {
            id: 'c4',
            name: 'Diana Prince',
            role: 'Product Manager',
            score: 88,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Diana',
            email: 'diana.p@example.com',
            appliedDate: '2024-02-14',
            matchSummary: 'Great leadership potential and product vision.',
        },
        {
            id: 'c5',
            name: 'Evan Wright',
            role: 'Designer',
            score: 78,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Evan',
            email: 'evan.w@example.com',
            appliedDate: '2024-02-11',
            matchSummary: 'Solid design fundamentals, needs more mobile experience.',
        },
    ],
    interview: [
        {
            id: 'c6',
            name: 'Fiona Gallagher',
            role: 'HR Specialist',
            score: 92,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Fiona',
            email: 'fiona.g@example.com',
            appliedDate: '2024-02-09',
            matchSummary: 'Excellent communication and empathy.',
        },
    ],
    offer: [
        {
            id: 'c7',
            name: 'George Martin',
            role: 'Tech Lead',
            score: 95,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=George',
            email: 'george.m@example.com',
            appliedDate: '2024-02-05',
            matchSummary: 'Perfect match for Tech Lead role.',
        },
    ],
    rejected: [
        {
            id: 'c8',
            name: 'Hannah Abbott',
            role: 'Intern',
            score: 45,
            avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Hannah',
            email: 'hannah.a@example.com',
            appliedDate: '2024-02-15',
            matchSummary: 'Too junior for current openings.',
        },
    ],
};
