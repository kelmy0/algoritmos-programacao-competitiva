import { API_URL } from "$env/static/private";
import { normalizeApiError } from "$lib/utils/errors";
import { error } from "@sveltejs/kit";
import { renderMarkdown } from "$lib/services/markdown";
import { customFetch } from "$lib/api/client";
import type { Algorithm } from "$lib/types/algorithm";
import type { PageServerLoad } from "./$types";
import { ALGORITHMS_ERRORS } from "$lib/errors/algorithms/algorithms";

export const load: PageServerLoad = async (event) => {
	const { slug } = event.params;

	const {
		data: responseBody,
		error: apiError,
		status,
		headers
	} = await customFetch<{ data: Algorithm }>(
		event.fetch,
		`${API_URL}/api/algorithms/${slug}`,
		{},
		ALGORITHMS_ERRORS
	);

	if (apiError) {
		error(status, apiError);
	}

	if (!responseBody || !responseBody.data) {
		error(500, normalizeApiError("INTERNAL_SERVER_ERROR"));
	}

	const cacheControl = headers.get("cache-control");
	if (cacheControl) {
		event.setHeaders({
			"cache-control": cacheControl
		});
	}

	const algorithmData = responseBody.data;
	const contentHtml = await renderMarkdown(algorithmData.content || "");

	return {
		algorithm: {
			...algorithmData,
			contentHtml
		}
	};
};
