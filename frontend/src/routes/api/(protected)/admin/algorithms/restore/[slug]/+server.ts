import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { normalizeApiError } from "$lib/utils/errors";
import { checkAdminAccess } from "$lib/utils/permissions";
import { json, type RequestHandler } from "@sveltejs/kit";

export const PATCH: RequestHandler = async (event) => {
	checkAdminAccess(event.locals.user, "create:algorithms");

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
		`${API_URL}/api/admin/algorithms/restore/${slug}`,
		{
			method: "PATCH",
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
