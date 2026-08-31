import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { ALGORITHMS_ERRORS } from "$lib/errors/algorithms/algorithms";
import { svelteServerCache } from "$lib/server/cache";
import { standardApiLimiter, thousandQuerySize, useMiddlewares } from "$lib/server/middlewares";
import { renderMarkdown } from "$lib/services/markdown";
import type { Algorithm, AlgorithmDetailResponse } from "$lib/types/algorithm";
import { normalizeApiError } from "$lib/utils/errors";
import { json, type RequestHandler } from "@sveltejs/kit";

const getAlgorithmBySlug: RequestHandler = async (event) => {
	const { slug } = event.params;

	const cacheKey = `algorithm:detail:${slug}`;
	const cachedResponse = svelteServerCache.get<AlgorithmDetailResponse>(cacheKey);
	if (cachedResponse) {
		return json(cachedResponse.data, {
			status: 200,
			headers: cachedResponse.headers
		});
	}

	const clientIp = event.getClientAddress();

	const {
		data: responseBody,
		error: ApiError,
		status,
		headers
	} = await customFetch<Algorithm>(
		event.fetch,
		`${API_URL}/api/algorithms/${slug}`,
		{
			headers: {
				"X-Forwarded-For": clientIp
			}
		},
		ALGORITHMS_ERRORS
	);

	if (ApiError) {
		return json(ApiError, { status });
	}

	if (!responseBody || !responseBody) {
		return json(normalizeApiError("ALGORITHM_NOT_FOUND"), { status: 404 });
	}

	const algorithmData = responseBody;
	const contentHtml = await renderMarkdown(algorithmData.content || "");

	const payload: AlgorithmDetailResponse = {
		algorithm: {
			...algorithmData,
			contentHtml
		}
	};

	const responseHeaders: Record<string, string> = {};
	const cacheControl = headers.get("cache-control");
	if (cacheControl) {
		responseHeaders["cache-control"] = cacheControl;
	}

	svelteServerCache.set(cacheKey, payload, responseHeaders, 30);

	return json(payload, { status, headers: responseHeaders });
};

export const GET = useMiddlewares(thousandQuerySize, standardApiLimiter)(getAlgorithmBySlug);
