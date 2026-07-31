import { algorithmSchema } from "$lib/schemas/algorithm";
import { checkAdminAccess } from "$lib/utils/permissions";
import { json, type RequestHandler } from "@sveltejs/kit";
import { ADMIN_ALGORITHMS_ERRORS } from "../../../../../(protected)/admin/algorithms/new/newAlgorithm.svelte";
import { normalizeApiError } from "$lib/utils/errors";
import { API_URL } from "$env/static/private";
import type { Algorithm } from "$lib/types/algorithm";
import { customFetch } from "$lib/api/client";

interface ApiResponse {
	data: Algorithm;
}

export const POST: RequestHandler = async (event) => {
	checkAdminAccess(event.locals.user, "create:algorithms");

	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_COOKIE",
			"Falta a senha das rotas admin.",
			ADMIN_ALGORITHMS_ERRORS
		);
		return json(normalizedError, { status: 401 });
	}

	const body = await event.request.json().catch(() => null);
	const result = algorithmSchema.safeParse(body);

	if (!result.success) {
		const errorCode = result.error.issues[0].message;
		const normalizedError = normalizeApiError(
			errorCode,
			"Formato de dados inválido.",
			ADMIN_ALGORITHMS_ERRORS
		);
		return json(normalizedError, { status: 400 });
	}

	const {
		data,
		error: apiError,
		status
	} = await customFetch<ApiResponse>(
		event.fetch,
		`${API_URL}/api/admin/algorithms`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Admin-Secret": adminSecret,
				Authorization: `Bearer ${event.locals.accessToken}`
			},
			body: JSON.stringify(result.data)
		},
		ADMIN_ALGORITHMS_ERRORS
	);

	if (apiError) {
		return json(apiError, { status });
	}

	return json({ algorithm: data?.data });
};
