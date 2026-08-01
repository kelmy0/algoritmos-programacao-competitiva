import { normalizeApiError } from "$lib/utils/errors";
import { type RequestHandler, json } from "@sveltejs/kit";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import { requireAuth, requirePermission, useMiddlewares } from "$lib/server/middlewares";

const deleteMyAlgorithm: RequestHandler = async (event) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_COOKIE",
			"Falta a senha das rotas admin."
		);
		return json(normalizedError, { status: 401 });
	}

	const { slug } = event.params;

	const { error: apiError, status } = await customFetch<null>(
		event.fetch,
		`${API_URL}/api/admin/algorithms/${slug}`,
		{
			method: "DELETE",
			headers: {
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

export const DELETE = useMiddlewares(
	requireAuth,
	requirePermission("create:algorithms")
)(deleteMyAlgorithm);
