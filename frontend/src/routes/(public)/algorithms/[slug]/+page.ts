import { PUBLIC_API_URL } from "$env/static/public";
import { normalizeApiError } from "$lib/utils/errors";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { ALGORITHMS_ERRORS } from "../algorithms";
import { renderMarkdown } from "$lib/services/markdown";
import { customFetch } from "$lib/api/client";
import type { Algorithm } from "$lib/types/algorithm";

export const load: PageLoad = async ({ fetch: svelteFetch, params }) => {
	const { slug } = params;

	const {
		data: responseBody,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(
		svelteFetch,
		`${PUBLIC_API_URL}/api/algorithms/${slug}`,
		{},
		ALGORITHMS_ERRORS
	);

	if (apiError) {
		error(status, apiError);
	}

	if (!responseBody || !responseBody.data) {
		error(
			500,
			normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Falha ao processar resposta do servidor.",
				ALGORITHMS_ERRORS
			)
		);
	}

	const algorithmData = responseBody.data;
	const contentHtml = await renderMarkdown(algorithmData.Content || "");

	return {
		algorithm: {
			...algorithmData,
			contentHtml
		}
	};
};
