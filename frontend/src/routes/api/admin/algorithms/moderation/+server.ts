import { normalizeApiError } from "$lib/utils/errors";
import { json, type RequestHandler } from "@sveltejs/kit";
import { rateLimit, requireAuth, requirePermission, useMiddlewares } from "$lib/server/middlewares";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";

interface ApiResponse {
	algorithms: Algorithm[];
	page: number;
	limit: number;
}

const moderationAlgorithms: RequestHandler = async (event) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_COOKIE",
			"Falta a senha das rotas admin."
		);
		return json(normalizedError, { status: 401 });
	}

	const queryString = event.url.searchParams.get("status");
	const clientIp = event.getClientAddress();

	const {
		data,
		error: apiError,
		status
	} = await customFetch<ApiResponse>(
		event.fetch,
		`${API_URL}/api/admin/algorithms/moderation${queryString ? `?status=${queryString}` : ""}`,
		{
			method: "GET",
			headers: {
				"X-Forwarded-For": clientIp,
				"X-Admin-Secret": adminSecret,
				Authorization: `Bearer ${event.locals.accessToken}`
			}
		}
	);

	if (apiError) {
		return json(apiError, { status });
	}

	if (!data?.algorithms) {
		return json({});
	}

	return json({ algorithms: data.algorithms, page: data.page, limit: data.limit });
};

export const GET = useMiddlewares(
	rateLimit({ capacity: 5, fillRate: 0.1 }),
	requireAuth,
	requirePermission("moderate:algorithms")
)(moderationAlgorithms);
