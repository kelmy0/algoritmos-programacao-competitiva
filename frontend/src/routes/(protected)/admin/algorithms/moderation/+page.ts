import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import type { Algorithm } from "$lib/types/algorithm";
import { customFetch } from "$lib/api/client";

interface ApiResponse {
	page: number;
	limit: number;
	algorithms: Algorithm[];
}

export const load: PageLoad = async (event) => {
	const {
		data,
		error: apiError,
		status
	} = await customFetch<ApiResponse>(event.fetch, `/api/admin/algorithms/moderation`);

	if (status === 401) {
		return {
			algorithms: [],
			pagination: { page: 1, limit: 10 }
		};
	}

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
			limit: data.limit
		}
	};
};
