import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { rateLimit, requireAuth, requirePermission, useMiddlewares } from "$lib/server/middlewares";
import { normalizeApiError } from "$lib/utils/errors";
import { json, type RequestHandler } from "@sveltejs/kit";

const restoreAlgorithm: RequestHandler = async (event) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_COOKIE",
			"Falta a senha das rotas admin."
		);
		return json(normalizedError, { status: 401 });
	}

	const { slug } = event.params;
	const clientIp = event.getClientAddress();

	const { error: apiError, status } = await customFetch<null>(
		event.fetch,
		`${API_URL}/api/admin/algorithms/restore/${slug}`,
		{
			method: "PATCH",
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

	return new Response(null, { status: 204 });
};

export const PATCH = useMiddlewares(
	rateLimit({ capacity: 5, fillRate: 0.1 }),
	requireAuth,
	requirePermission("create:algorithms")
)(restoreAlgorithm);
