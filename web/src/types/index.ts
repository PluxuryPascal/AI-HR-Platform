// Global TypeScript interfaces and types

/**
 * Standard API response envelope from the backend.
 */
export interface ApiResponse<T> {
    data: T;
    message?: string;
    success: boolean;
}

/**
 * Pagination metadata returned by list endpoints.
 */
export interface PaginationMeta {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
}

/**
 * Paginated API response.
 */
export interface PaginatedResponse<T> extends ApiResponse<T[]> {
    meta: PaginationMeta;
}
