import { normalizeApiError } from "$lib/utils/errors";
import { error } from "@sveltejs/kit";
import { customFetch } from "$lib/api/client";
import type { AlgorithmDetailResponse } from "$lib/types/algorithm";
import type { PageLoad } from "./$types";
import { ALGORITHMS_ERRORS } from "$lib/errors/algorithms/algorithms";

export const load: PageLoad = async (event) => {
	const { slug } = event.params;

	const {
		data,
		error: apiError,
		status
	} = await customFetch<AlgorithmDetailResponse>(
		event.fetch,
		`/api/algorithms/${slug}`,
		{},
		ALGORITHMS_ERRORS
	);

	if (apiError) {
		error(status, apiError);
	}

	if (!data) {
		error(500, normalizeApiError("INTERNAL_SERVER_ERROR"));
	}

	return {
		algorithm: data.algorithm
	};
};
