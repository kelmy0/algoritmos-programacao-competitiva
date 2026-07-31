import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import type { Algorithm } from "$lib/types/algorithm";
import { normalizeApiError } from "$lib/utils/errors";
import { checkAdminAccess } from "$lib/utils/permissions";
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

interface ApiResponse {
	algorithms: Algorithm[];
	page: number;
	limit: number;
}

export const GET: RequestHandler = async (event) => {
	checkAdminAccess(event.locals.user, "create:algorithms");

	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_COOKIE",
			"Falta a senha das rotas admin."
		);
		return json(normalizedError, { status: 401 });
	}

	const queryString = event.url.searchParams.get("status");

	const {
		data,
		error: apiError,
		status
	} = await customFetch<ApiResponse>(
		event.fetch,
		`${API_URL}/api/admin/algorithms${queryString ? `?status=${queryString}` : ""}`,
		{
			method: "GET",
			headers: {
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
