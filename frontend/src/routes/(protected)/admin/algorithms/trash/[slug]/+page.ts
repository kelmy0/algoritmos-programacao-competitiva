import { customFetch } from "$lib/api/client";
import type { Algorithm } from "$lib/types/algorithm";
import { error } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import type { PageLoad } from "./$types";
import { renderMarkdown } from "$lib/services/markdown";

export const load: PageLoad = async ({ fetch: svelteFetch, params }) => {
	const { slug } = params;

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(svelteFetch, `/api/admin/algorithms/edit/${slug}`);

	if (status === 401 || apiError?.code === "MISSING_ADMIN_COOKIE") {
		return { algorithm: null };
	}

	if (apiError) {
		error(status, apiError);
	}

	if (!data || !data.data) {
		error(
			500,
			normalizeApiError("INTERNAL_SERVER_ERROR", "Falha ao processar resposta do servidor.")
		);
	}

	const algorithmData = data.data;
	const contentHtml = await renderMarkdown(algorithmData.Content || "");

	return {
		algorithm: {
			...algorithmData,
			contentHtml
		}
	};
};
