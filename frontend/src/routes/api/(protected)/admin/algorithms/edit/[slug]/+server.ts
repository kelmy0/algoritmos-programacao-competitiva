import { normalizeApiError } from "$lib/utils/errors";
import { checkAdminAccess } from "$lib/utils/permissions";
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { customFetch } from "$lib/api/client";
import { PUBLIC_API_URL } from "$env/static/public";
import type { Algorithm } from "$lib/types/algorithm";

export const GET: RequestHandler = async (event) => {
	checkAdminAccess(event.locals.user, "create:algorithms");

	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError("MISSING_COOKIE", "Falta a senha das rotas admin.");
		return json(normalizedError, { status: 400 });
	}

	const { slug } = event.params;

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(
		event.fetch,
		`${PUBLIC_API_URL}/api/admin/algorithms/${slug}`,
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

	if (!data?.data) {
		return json({});
	}

	return json(data);
};
