import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import type { ListAlgorithmsResponse } from "$lib/types/algorithm";
import { customFetch } from "$lib/api/client";
import { ALGORITHMS_ERRORS } from "$lib/errors/algorithms/algorithms";

export const load: PageLoad = async (event) => {
	const {
		data,
		error: apiError,
		status
	} = await customFetch<ListAlgorithmsResponse>(
		event.fetch,
		`/api/algorithms`,
		{},
		ALGORITHMS_ERRORS
	);

	if (apiError) {
		error(status, apiError);
	}

	if (!data) {
		return {};
	}

	return {
		algorithms: data.algorithms,
		pagination: {
			page: data.page,
			limit: data.limit,
			hasMore: data.hasMore
		}
	};
};
