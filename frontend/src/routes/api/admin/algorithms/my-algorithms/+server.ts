import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import type { ListAlgorithmsResponse } from "$lib/types/algorithm";
import { normalizeApiError } from "$lib/utils/errors";
import { json, type RequestHandler } from "@sveltejs/kit";
import { standardApiLimiter, requireAuth, requirePermission } from "$lib/server/middlewares";
import { thousandQuerySize, useMiddlewares } from "$lib/server/middlewares";
import { ADMIN_PASSWORD_ERRORS } from "$lib/errors/admin/password";

const myAlgorithms: RequestHandler = async (event) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError("MISSING_ADMIN_COOKIE", "", ADMIN_PASSWORD_ERRORS);
		return json(normalizedError, { status: 401 });
	}

	const queryString = event.url.searchParams.get("status");
	const clientIp = event.getClientAddress();

	const {
		data,
		error: apiError,
		status
	} = await customFetch<ListAlgorithmsResponse>(
		event.fetch,
		`${API_URL}/api/admin/algorithms${queryString ? `?status=${queryString}` : ""}`,
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
	thousandQuerySize,
	standardApiLimiter,
	requireAuth,
	requirePermission("create:algorithms")
)(myAlgorithms);
