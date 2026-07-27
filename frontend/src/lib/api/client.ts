import { browser } from "$app/environment";
import { goto, invalidateAll } from "$app/navigation";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";

export async function customFetch<T>(
	fetchImpl: typeof fetch,
	url: string,
	options?: RequestInit,
	localErrors?: Record<string, string>
): Promise<{ data: T | null; error: ApiError | null; status: number; headers: Headers }> {
	try {
		const response = await fetchImpl(url, options);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({}));
			const normalizedError = normalizeApiError(
				errorData,
				"Ocorreu um erro no servidor.",
				localErrors
			);

			if (normalizedError.code === "TOKEN_NO_LONGER_VALID") {
				if (browser) {
					await goto("/auth/login?error=TOKEN_NO_LONGER_VALID");
					await invalidateAll();
				}
			}

			return {
				data: null,
				error: normalizedError,
				status: response.status,
				headers: response.headers
			};
		}

		if (response.status === 204) {
			return { data: null, error: null, status: 204, headers: response.headers };
		}

		const data: T = await response.json().catch(() => null as unknown as T);
		return { data, error: null, status: response.status, headers: response.headers };
	} catch (err) {
		return {
			data: null,
			error: normalizeApiError(err, "Falha ao se comunicar com o servidor.", localErrors),
			status: 500,
			headers: new Headers()
		};
	}
}
