/**
 * Centralised API client for communicating with backend micro-services.
 *
 * Every feature's `api/` layer should import from here rather than
 * calling `fetch` directly, so that base URL, headers, and error
 * handling stay consistent across the app.
 */

const API_BASE_URL =
    process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api/v1";

interface RequestOptions extends Omit<RequestInit, "body"> {
    body?: unknown;
}

async function request<T>(
    endpoint: string,
    options: RequestOptions = {},
): Promise<T> {
    const { body, headers, ...rest } = options;

    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
        headers: {
            "Content-Type": "application/json",
            ...headers,
        },
        body: body ? JSON.stringify(body) : undefined,
        ...rest,
    });

    if (!res.ok) {
        if (res.status === 401) {
            if (typeof window !== "undefined") {
                const pathname = window.location.pathname;
                if (!pathname.includes("/auth")) {
                    const match = pathname.match(/^\/(en|ru|es)(\/|$)/);
                    const localePrefix = match ? `/${match[1]}` : "";
                    window.location.href = `${localePrefix}/auth`;
                }
            }
        }

        const error = await res.json().catch(() => ({}));
        throw new Error(
            (error as { message?: string }).message ??
            `API Error: ${res.status} ${res.statusText}`,
        );
    }

    const text = await res.text();
    if (!text) return undefined as T;

    try {
        return JSON.parse(text) as T;
    } catch (e) {
        return text as unknown as T;
    }
}

export const apiClient = {
    get: <T>(endpoint: string, opts?: RequestOptions) =>
        request<T>(endpoint, { ...opts, method: "GET" }),

    post: <T>(endpoint: string, body?: unknown, opts?: RequestOptions) =>
        request<T>(endpoint, { ...opts, method: "POST", body }),

    put: <T>(endpoint: string, body?: unknown, opts?: RequestOptions) =>
        request<T>(endpoint, { ...opts, method: "PUT", body }),

    patch: <T>(endpoint: string, body?: unknown, opts?: RequestOptions) =>
        request<T>(endpoint, { ...opts, method: "PATCH", body }),

    delete: <T>(endpoint: string, opts?: RequestOptions) =>
        request<T>(endpoint, { ...opts, method: "DELETE" }),
};
