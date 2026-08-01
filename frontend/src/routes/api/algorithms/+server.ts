import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { standardApiLimiter, thousandQuerySize, useMiddlewares } from "$lib/server/middlewares";
import type { ListAlgorithmsResponse } from "$lib/types/algorithm";
import { json, type RequestHandler } from "@sveltejs/kit";

const listAlgorithms: RequestHandler = async (event) => {
	const rawPage = parseInt(event.url.searchParams.get("page") ?? "1", 10);
	const rawLimit = parseInt(event.url.searchParams.get("limit") ?? "12", 10);

	const page = Math.max(1, isNaN(rawPage) ? 1 : rawPage);
	const limit = Math.min(50, Math.max(1, isNaN(rawLimit) ? 10 : rawLimit));

	const clientIp = event.getClientAddress();

	const {
		data,
		error: ApiError,
		status,
		headers
	} = await customFetch<ListAlgorithmsResponse>(
		event.fetch,
		`${API_URL}/api/algorithms?page=${page}&limit=${limit}`,
		{
			headers: {
				"X-Forwarded-For": clientIp
			}
		}
	);

	if (ApiError) {
		return json(ApiError, { status });
	}

	const responseHeaders: Record<string, string> = {};
	const cacheControl = headers.get("cache-control");
	if (cacheControl) {
		responseHeaders["cache-control"] = cacheControl;
	}

	return json(data, { status, headers });
};

export const GET = useMiddlewares(thousandQuerySize, standardApiLimiter)(listAlgorithms);
