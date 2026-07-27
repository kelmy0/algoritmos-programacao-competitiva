import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";

export async function customFetch<T>(
	fetchImpl: typeof fetch,
	url: string,
	options?: RequestInit,
	localErrors?: Record<string, string>
): Promise<{ data: T | null; error: ApiError | null; status: number }> {
	try {
		const response = await fetchImpl(url, options);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({}));

			return {
				data: null,
				error: normalizeApiError(errorData, "Ocorreu um erro no servidor.", localErrors),
				status: response.status
			};
		}

		if (response.status === 204) {
			return { data: null, error: null, status: 204 };
		}

		const data: T = await response.json().catch(() => null as unknown as T);
		return { data, error: null, status: response.status };
	} catch (err) {
		return {
			data: null,
			error: normalizeApiError(err, "Falha ao se comunicar com o servidor.", localErrors),
			status: 500
		};
	}
}
