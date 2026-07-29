import { customFetch } from "$lib/api/client";
import type { Algorithm } from "$lib/types/algorithm";
import { error } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch: svelteFetch, params }) => {
	const { slug } = params;

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(svelteFetch, `/api/admin/algorithms/edit/${slug}`);

	if (apiError) {
		error(status, apiError);
	}

	if (!data || !data.data) {
		error(
			500,
			normalizeApiError("INTERNAL_SERVER_ERROR", "Falha ao processar resposta do servidor.")
		);
	}

	return data.data;
};
