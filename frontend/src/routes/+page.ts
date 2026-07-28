import { PUBLIC_API_URL } from "$env/static/public";
import { normalizeApiError } from "$lib/utils/errors";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import type { Algorithm } from "$lib/types/algorithm";
import { customFetch } from "$lib/api/client";

interface ApiResponse {
	page: number;
	limit: number;
	data: Algorithm[];
}

export async function load({ fetch: svelteFetch }: Parameters<PageLoad>[0]) {
	const {
		data,
		error: apiError,
		status
	} = await customFetch<ApiResponse>(svelteFetch, `${PUBLIC_API_URL}/api/algorithms`, {});

	if (apiError) {
		error(status, apiError);
	}

	if (!data || !data.data) {
		error(
			500,
			normalizeApiError("INTERNAL_SERVER_ERROR", "Falha ao processar resposta do servidor.")
		);
	}

	return {
		algorithms: data.data,
		pagination: {
			page: data.page,
			limit: data.limit
		}
	};
}
