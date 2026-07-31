import { PUBLIC_API_URL } from "$env/static/public";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import type { Algorithm } from "$lib/types/algorithm";
import { customFetch } from "$lib/api/client";

interface ApiResponse {
	page: number;
	limit: number;
	data: Algorithm[];
}

export const load: PageLoad = async ({ fetch: svelteFetch }) => {
	const {
		data,
		error: apiError,
		status
	} = await customFetch<ApiResponse>(svelteFetch, `${PUBLIC_API_URL}/api/algorithms`, {});

	if (apiError) {
		error(status, apiError);
	}

	if (!data) {
		return {};
	}

	return {
		algorithms: data.data,
		pagination: {
			page: data.page,
			limit: data.limit
		}
	};
};
