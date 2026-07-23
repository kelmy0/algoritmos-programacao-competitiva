import type { ApiError } from '$lib/types/api';
import { normalizeApiError } from '$lib/utils/errors';

export async function customFetch<T>(
	fetchImpl: typeof fetch,
	url: string,
	options?: RequestInit,
	localErrors?: Record<string, string>
): Promise<{ data: T | null; error: ApiError | null }> {
	try {
		const response = await fetchImpl(url, options);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({}));

			return {
				data: null,
				error: normalizeApiError(errorData, 'Ocorreu um erro no servidor.', localErrors)
			};
		}

		const data = await response.json();
		return { data, error: null };
	} catch (err) {
		return {
			data: null,
			error: normalizeApiError(err, 'Falha ao se comunicar com o servidor.', localErrors)
		};
	}
}
