import { customFetch } from "$lib/api/client";
import type { Algorithm } from "$lib/types/algorithm";
import { error } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import type { PageLoad } from "./$types";

export const load: PageLoad = async (event) => {
	const { slug } = event.params;

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(event.fetch, `/api/admin/algorithms/edit/${slug}`);

	if (status === 401 || apiError?.code === "MISSING_ADMIN_COOKIE") {
		return { algorithm: null };
	}

	if (apiError) {
		error(status, apiError);
	}

	if (!data || !data.data) {
		error(500, normalizeApiError("INTERNAL_SERVER_ERROR"));
	}

	return { algorithm: data.data };
};
